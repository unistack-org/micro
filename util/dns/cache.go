package dns

import (
	"context"
	"math"
	"net"
	"sync"
	"time"

	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/semconv"
)

// DialFunc is a [net.Resolver.Dial] function.
type DialFunc func(ctx context.Context, network, address string) (net.Conn, error)

// NewNetResolver creates a caching [net.Resolver] that uses parent to resolve names.
func NewNetResolver(opts ...Option) *net.Resolver {
	options := Options{Resolver: &net.Resolver{}}

	for _, o := range opts {
		o(&options)
	}

	if options.Meter == nil {
		options.Meter = meter.DefaultMeter
		opts = append(opts, Meter(options.Meter))
	}

	return &net.Resolver{
		PreferGo:     true,
		StrictErrors: options.Resolver.StrictErrors,
		Dial:         NewNetDialer(options.Resolver.Dial, append(opts, Resolver(options.Resolver))...),
	}
}

// NewNetDialer adds caching to a [net.Resolver.Dial] function.
func NewNetDialer(parent DialFunc, opts ...Option) DialFunc {
	cache := cache{dial: parent, opts: Options{}}
	for _, o := range opts {
		o(&cache.opts)
	}
	if cache.opts.MaxCacheEntries == 0 {
		cache.opts.MaxCacheEntries = DefaultMaxCacheEntries
	}
	return func(_ context.Context, network, address string) (net.Conn, error) {
		conn := &dnsConn{}
		conn.roundTrip = cachingRoundTrip(&cache, network, address)
		return conn, nil
	}
}

const DefaultMaxCacheEntries = 300

// A Option customizes the resolver cache.
type Option func(*Options)

type Options struct {
	Resolver        *net.Resolver
	MaxCacheEntries int
	MaxCacheTTL     time.Duration
	MinCacheTTL     time.Duration
	NegativeCache   bool
	PreferIPV4      bool
	PreferIPV6      bool
	Timeout         time.Duration
	Meter           meter.Meter
}

// MaxCacheEntries sets the maximum number of entries to cache.
// If zero, [DefaultMaxCacheEntries] is used; negative means no limit.
func MaxCacheEntries(n int) Option {
	return func(o *Options) {
		o.MaxCacheEntries = n
	}
}

// MaxCacheTTL sets the maximum time-to-live for entries in the cache.
func MaxCacheTTL(td time.Duration) Option {
	return func(o *Options) {
		o.MaxCacheTTL = td
	}
}

// MinCacheTTL sets the minimum time-to-live for entries in the cache.
func MinCacheTTL(td time.Duration) Option {
	return func(o *Options) {
		o.MinCacheTTL = td
	}
}

// NegativeCache sets whether to cache negative responses.
func NegativeCache(b bool) Option {
	return func(o *Options) {
		o.NegativeCache = b
	}
}

// Meter sets meter.Meter
func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// Timeout sets upstream *net.Resolver timeout
func Timeout(td time.Duration) Option {
	return func(o *Options) {
		o.Timeout = td
	}
}

// Resolver sets upstream *net.Resolver.
func Resolver(r *net.Resolver) Option {
	return func(o *Options) {
		o.Resolver = r
	}
}

// PreferIPV4 resolve ipv4 records.
func PreferIPV4(b bool) Option {
	return func(o *Options) {
		o.PreferIPV4 = b
	}
}

// PreferIPV6 resolve ipv4 records.
func PreferIPV6(b bool) Option {
	return func(o *Options) {
		o.PreferIPV6 = b
	}
}

type cache struct {
	entries map[string]cacheEntry
	dial    DialFunc

	opts Options

	mu sync.RWMutex
}

type cacheEntry struct {
	deadline time.Time
	value    string
}

func (c *cache) put(req string, res string) {
	// ignore uncacheable/unparseable answers
	if invalid(req, res) {
		return
	}

	// ignore errors (if requested)
	if nameError(res) && !c.opts.NegativeCache {
		return
	}

	// ignore uncacheable/unparseable answers
	ttl := getTTL(res)
	if ttl <= 0 {
		return
	}

	// adjust TTL
	if ttl < c.opts.MinCacheTTL {
		ttl = c.opts.MinCacheTTL
	}
	// maxTTL overrides minTTL
	if ttl > c.opts.MaxCacheTTL && c.opts.MaxCacheTTL != 0 {
		ttl = c.opts.MaxCacheTTL
	}

	c.mu.Lock()
	if c.entries == nil {
		c.entries = make(map[string]cacheEntry)
	}

	// do some cache evition
	var tested, evicted int
	for k, e := range c.entries {
		if time.Until(e.deadline) <= 0 {
			c.opts.Meter.Counter(semconv.CacheItemsTotal, "type", "dns").Dec()
			c.opts.Meter.Counter(semconv.CacheRequestTotal, "type", "dns", "method", "evict").Inc()
			// delete expired entry
			delete(c.entries, k)
			evicted++
		}
		tested++

		if tested < 8 {
			continue
		}
		if evicted == 0 && c.opts.MaxCacheEntries > 0 && len(c.entries) >= c.opts.MaxCacheEntries {
			c.opts.Meter.Counter(semconv.CacheItemsTotal, "type", "dns").Dec()
			c.opts.Meter.Counter(semconv.CacheRequestTotal, "type", "dns", "method", "evict").Inc()
			// delete at least one entry
			delete(c.entries, k)
		}
		break
	}

	// remove message IDs
	c.entries[req[2:]] = cacheEntry{
		deadline: time.Now().Add(ttl),
		value:    res[2:],
	}

	c.opts.Meter.Counter(semconv.CacheItemsTotal, "type", "dns").Inc()
	c.mu.Unlock()
}

func (c *cache) get(req string) (res string) {
	// ignore invalid messages
	if len(req) < 12 {
		return ""
	}
	if req[2] >= 0x7f {
		return ""
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.entries == nil {
		return ""
	}

	// remove message ID
	entry, ok := c.entries[req[2:]]
	if ok && time.Until(entry.deadline) > 0 {
		// prepend correct ID
		return req[:2] + entry.value
	}

	return ""
}

func invalid(req string, res string) bool {
	if len(req) < 12 || len(res) < 12 { // header size
		return true
	}
	if req[0] != res[0] || req[1] != res[1] { // IDs match
		return true
	}
	if req[2] >= 0x7f || res[2] < 0x7f { // query, response
		return true
	}
	if req[2]&0x7a != 0 || res[2]&0x7a != 0 { // standard query, not truncated
		return true
	}
	if res[3]&0xf != 0 && res[3]&0xf != 3 { // no error, or name error
		return true
	}
	return false
}

func nameError(res string) bool {
	return res[3]&0xf == 3
}

func getTTL(msg string) time.Duration {
	ttl := math.MaxInt32

	qdcount := getUint16(msg[4:])
	ancount := getUint16(msg[6:])
	nscount := getUint16(msg[8:])
	arcount := getUint16(msg[10:])
	rdcount := ancount + nscount + arcount

	msg = msg[12:] // skip header

	// skip questions
	for i := 0; i < qdcount; i++ {
		name := getNameLen(msg)
		if name < 0 || name+4 > len(msg) {
			return -1
		}
		msg = msg[name+4:]
	}

	// parse records
	for i := 0; i < rdcount; i++ {
		name := getNameLen(msg)
		if name < 0 || name+10 > len(msg) {
			return -1
		}
		rtyp := getUint16(msg[name+0:])
		rttl := getUint32(msg[name+4:])
		rlen := getUint16(msg[name+8:])
		if name+10+rlen > len(msg) {
			return -1
		}
		// skip EDNS OPT since it doesn't have a TTL
		if rtyp != 41 && rttl < ttl {
			ttl = rttl
		}
		msg = msg[name+10+rlen:]
	}

	return time.Duration(ttl) * time.Second
}

func getNameLen(msg string) int {
	i := 0
	for i < len(msg) {
		if msg[i] == 0 {
			// end of name
			i++
			break
		}
		if msg[i] >= 0xc0 {
			// compressed name
			i += 2
			break
		}
		if msg[i] >= 0x40 {
			// reserved
			return -1
		}
		i += int(msg[i] + 1)
	}
	return i
}

func getUint16(s string) int {
	return int(s[1]) | int(s[0])<<8
}

func getUint32(s string) int {
	return int(s[3]) | int(s[2])<<8 | int(s[1])<<16 | int(s[0])<<24
}

func cachingRoundTrip(cache *cache, network, address string) roundTripper {
	return func(ctx context.Context, req string) (res string, err error) {
		cache.opts.Meter.Counter(semconv.CacheRequestInflight, "type", "dns").Inc()
		defer cache.opts.Meter.Counter(semconv.CacheRequestInflight, "type", "dns").Dec()
		// check cache
		if res = cache.get(req); res != "" {
			return res, nil
		}
		cache.opts.Meter.Counter(semconv.CacheRequestTotal, "type", "dns", "method", "get", "status", "miss").Inc()
		ts := time.Now()
		defer func() {
			cache.opts.Meter.Summary(semconv.CacheRequestLatencyMicroseconds, "type", "dns", "method", "get").UpdateDuration(ts)
			cache.opts.Meter.Histogram(semconv.CacheRequestDurationSeconds, "type", "dns", "method", "get").UpdateDuration(ts)
		}()

		switch {
		case cache.opts.PreferIPV4 && cache.opts.PreferIPV6:
			network = "udp"
		case cache.opts.PreferIPV4:
			network = "udp4"
		case cache.opts.PreferIPV6:
			network = "udp6"
		default:
			network = "udp"
		}

		if cache.opts.Timeout > 0 {
			var cancel func()
			ctx, cancel = context.WithTimeout(ctx, cache.opts.Timeout)
			defer cancel()
		}

		// dial connection
		var conn net.Conn
		if cache.dial != nil {
			conn, err = cache.dial(ctx, network, address)
		} else {
			var d net.Dialer
			conn, err = d.DialContext(ctx, network, address)
		}

		if err != nil {
			return "", err
		}

		ctx, cancel := context.WithCancel(ctx)
		go func() {
			<-ctx.Done()
			conn.Close()
		}()
		defer cancel()

		if t, ok := ctx.Deadline(); ok {
			err = conn.SetDeadline(t)
			if err != nil {
				return "", err
			}
		}

		// send request
		err = writeMessage(conn, req)
		if err != nil {
			return "", err
		}

		// read response
		res, err = readMessage(conn)
		if err != nil {
			return "", err
		}

		// cache response
		cache.put(req, res)
		return res, nil
	}
}
