package unwrap

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"go.unistack.org/micro/v3/codec"
)

const sf = "0-+# "

var hexDigits = "0123456789abcdef"

var (
	panicBytes         = []byte("(PANIC=")
	plusBytes          = []byte("+")
	iBytes             = []byte("i")
	trueBytes          = []byte("true")
	falseBytes         = []byte("false")
	interfaceBytes     = []byte("(interface {})")
	openBraceBytes     = []byte("{")
	closeBraceBytes    = []byte("}")
	asteriskBytes      = []byte("*")
	ampBytes           = []byte("&")
	colonBytes         = []byte(":")
	openParenBytes     = []byte("(")
	closeParenBytes    = []byte(")")
	spaceBytes         = []byte(" ")
	commaBytes         = []byte(",")
	pointerChainBytes  = []byte("->")
	nilAngleBytes      = []byte("<nil>")
	circularShortBytes = []byte("<shown>")
	invalidAngleBytes  = []byte("<invalid>")
	filteredBytes      = []byte("<filtered>")
	// openBracketBytes   = []byte("[")
	// closeBracketBytes  = []byte("]")
	percentBytes    = []byte("%")
	precisionBytes  = []byte(".")
	openAngleBytes  = []byte("<")
	closeAngleBytes = []byte(">")
	openMapBytes    = []byte("{")
	closeMapBytes   = []byte("}")
)

type protoMessage interface {
	Reset()
	ProtoMessage()
}

type Wrapper struct {
	val              interface{}
	s                fmt.State
	pointers         map[uintptr]int
	opts             *Options
	depth            int
	ignoreNextType   bool
	takeMap          map[int]bool
	protoWrapperType bool
	sqlWrapperType   bool
}

// Options struct
type Options struct {
	Codec   codec.Codec
	Indent  string
	Methods bool
	Tagged  bool
}

// NewOptions creates new Options struct via provided args
func NewOptions(opts ...Option) Options {
	options := Options{
		Indent:  " ",
		Methods: false,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Option func signature
type Option func(*Options)

// Indent option specify indent level
func Indent(f string) Option {
	return func(o *Options) {
		o.Indent = f
	}
}

// Methods option toggles fmt.Stringer methods
func Methods(b bool) Option {
	return func(o *Options) {
		o.Methods = b
	}
}

// Codec option automatic marshal arg via specified codec and write it to log
func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

// Tagged option toggles output only logger:"take" fields
func Tagged(b bool) Option {
	return func(o *Options) {
		o.Tagged = b
	}
}

func Unwrap(val interface{}, opts ...Option) *Wrapper {
	options := NewOptions(opts...)
	return &Wrapper{val: val, opts: &options, pointers: make(map[uintptr]int), takeMap: make(map[int]bool)}
}

func (w *Wrapper) unpackValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		w.ignoreNextType = false
		if !v.IsNil() {
			v = v.Elem()
		}
	}
	return v
}

// formatPtr handles formatting of pointers by indirecting them as necessary.
func (w *Wrapper) formatPtr(v reflect.Value) {
	// Display nil if top level pointer is nil.
	showTypes := w.s.Flag('#')
	if v.IsNil() && (!showTypes || w.ignoreNextType) {
		_, _ = w.s.Write(nilAngleBytes)
		return
	}

	// Remove pointers at or below the current depth from map used to detect
	// circular refs.
	for k, depth := range w.pointers {
		if depth >= w.depth {
			delete(w.pointers, k)
		}
	}

	// Keep list of all dereferenced pointers to possibly show later.
	pointerChain := make([]uintptr, 0)

	// Figure out how many levels of indirection there are by derferencing
	// pointers and unpacking interfaces down the chain while detecting circular
	// references.
	nilFound := false
	cycleFound := false
	indirects := 0
	ve := v
	for ve.Kind() == reflect.Ptr {
		if ve.IsNil() {
			nilFound = true
			break
		}
		indirects++
		addr := ve.Pointer()
		pointerChain = append(pointerChain, addr)
		if pd, ok := w.pointers[addr]; ok && pd < w.depth {
			cycleFound = true
			indirects--
			break
		}
		w.pointers[addr] = w.depth

		ve = ve.Elem()
		if ve.Kind() == reflect.Interface {
			if ve.IsNil() {
				nilFound = true
				break
			}
			ve = ve.Elem()
		}
	}

	// Display type or indirection level depending on flags.
	if showTypes && !w.ignoreNextType {
		if w.depth > 0 {
			_, _ = w.s.Write(openParenBytes)
		}
		if w.depth > 0 {
			_, _ = w.s.Write(bytes.Repeat(asteriskBytes, indirects))
		} else {
			_, _ = w.s.Write(bytes.Repeat(ampBytes, indirects))
		}
		_, _ = w.s.Write([]byte(ve.Type().String()))
		if w.depth > 0 {
			_, _ = w.s.Write(closeParenBytes)
		}
	} else {
		if nilFound || cycleFound {
			indirects += strings.Count(ve.Type().String(), "*")
		}
		_, _ = w.s.Write(openAngleBytes)
		_, _ = w.s.Write([]byte(strings.Repeat("*", indirects)))
		_, _ = w.s.Write(closeAngleBytes)
	}

	// Display pointer information depending on flags.
	if w.s.Flag('+') && (len(pointerChain) > 0) {
		_, _ = w.s.Write(openParenBytes)
		for i, addr := range pointerChain {
			if i > 0 {
				_, _ = w.s.Write(pointerChainBytes)
			}
			getHexPtr(w.s, addr)
		}
		_, _ = w.s.Write(closeParenBytes)
	}

	// Display dereferenced value.
	switch {
	case nilFound:
		_, _ = w.s.Write(nilAngleBytes)
	case cycleFound:
		_, _ = w.s.Write(circularShortBytes)
	default:
		w.ignoreNextType = true
		w.format(ve)
	}
}

// format is the main workhorse for providing the Formatter interface.  It
// uses the passed reflect value to figure out what kind of object we are
// dealing with and formats it appropriately.  It is a recursive function,
// however circular data structures are detected and handled properly.
func (w *Wrapper) format(v reflect.Value) {
	if w.opts.Codec != nil {
		buf, err := w.opts.Codec.Marshal(v.Interface())
		if err != nil {
			_, _ = w.s.Write(invalidAngleBytes)
			return
		}
		_, _ = w.s.Write(buf)
		return
	}

	// Handle invalid reflect values immediately.
	kind := v.Kind()
	if kind == reflect.Invalid {
		_, _ = w.s.Write(invalidAngleBytes)
		return
	}

	// Handle pointers specially.
	switch kind {
	case reflect.Ptr:
		if !v.IsZero() {
			if strings.HasPrefix(reflect.Indirect(v).Type().String(), "wrapperspb.") {
				w.protoWrapperType = true
			} else if strings.HasPrefix(reflect.Indirect(v).Type().String(), "sql.Null") {
				w.sqlWrapperType = true
			} else if v.CanInterface() {
				if _, ok := v.Interface().(protoMessage); ok {
					w.protoWrapperType = true
				}
			}
		}
		w.formatPtr(v)
		return
	case reflect.Struct:
		if !v.IsZero() {
			if strings.HasPrefix(reflect.Indirect(v).Type().String(), "sql.Null") {
				w.sqlWrapperType = true
			}
		}
	}

	// get type information unless already handled elsewhere.
	if !w.ignoreNextType && w.s.Flag('#') {
		if v.Type().Kind() != reflect.Map &&
			v.Type().Kind() != reflect.String &&
			v.Type().Kind() != reflect.Array &&
			v.Type().Kind() != reflect.Slice {
			_, _ = w.s.Write(openParenBytes)
		}
		if v.Kind() != reflect.String {
			_, _ = w.s.Write([]byte(v.Type().String()))
		}
		if v.Type().Kind() != reflect.Map &&
			v.Type().Kind() != reflect.String &&
			v.Type().Kind() != reflect.Array &&
			v.Type().Kind() != reflect.Slice {
			_, _ = w.s.Write(closeParenBytes)
		}
	}
	w.ignoreNextType = false

	// Call Stringer/error interfaces if they exist and the handle methods
	// flag is enabled.
	if w.opts.Methods {
		if (kind != reflect.Invalid) && (kind != reflect.Interface) {
			if handled := handleMethods(w.opts, w.s, v); handled {
				return
			}
		}
	}

	switch kind {
	case reflect.Invalid:
		_, _ = w.s.Write(invalidAngleBytes)
	case reflect.Bool:
		getBool(w.s, v.Bool())
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		getInt(w.s, v.Int(), 10)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		getUint(w.s, v.Uint(), 10)
	case reflect.Float32:
		getFloat(w.s, v.Float(), 32)
	case reflect.Float64:
		getFloat(w.s, v.Float(), 64)
	case reflect.Complex64:
		getComplex(w.s, v.Complex(), 32)
	case reflect.Complex128:
		getComplex(w.s, v.Complex(), 64)
	case reflect.Slice:
		if v.IsNil() {
			_, _ = w.s.Write(nilAngleBytes)
			break
		}
		fallthrough
	case reflect.Array:
		_, _ = w.s.Write(openBraceBytes)
		w.depth++
		numEntries := v.Len()
		for i := 0; i < numEntries; i++ {
			if i > 0 {
				_, _ = w.s.Write(commaBytes)
				_, _ = w.s.Write(spaceBytes)
			}
			w.ignoreNextType = true
			w.format(w.unpackValue(v.Index(i)))
		}
		w.depth--
		_, _ = w.s.Write(closeBraceBytes)
	case reflect.String:
		_, _ = w.s.Write([]byte(`"` + v.String() + `"`))
	case reflect.Interface:
		// The only time we should get here is for nil interfaces due to
		// unpackValue calls.
		if v.IsNil() {
			_, _ = w.s.Write(nilAngleBytes)
		}
	case reflect.Ptr:
		// Do nothing.  We should never get here since pointers have already
		// been handled above.
	case reflect.Map:
		// nil maps should be indicated as different than empty maps
		if v.IsNil() {
			_, _ = w.s.Write(nilAngleBytes)
			break
		}
		_, _ = w.s.Write(openMapBytes)
		w.depth++
		keys := v.MapKeys()
		for i, key := range keys {
			if i > 0 {
				_, _ = w.s.Write(spaceBytes)
			}
			w.ignoreNextType = true
			w.format(w.unpackValue(key))
			_, _ = w.s.Write(colonBytes)
			w.ignoreNextType = true
			w.format(w.unpackValue(v.MapIndex(key)))
		}
		w.depth--
		_, _ = w.s.Write(closeMapBytes)
	case reflect.Struct:

		numFields := v.NumField()
		numWritten := 0
		_, _ = w.s.Write(openBraceBytes)
		w.depth++

		vt := v.Type()
		prevSkip := false

		for i := 0; i < numFields; i++ {
			switch vt.Field(i).Type.PkgPath() {
			case "google.golang.org/protobuf/internal/impl", "google.golang.org/protobuf/internal/pragma":
				w.protoWrapperType = true
				prevSkip = true
				continue
			}
			if w.protoWrapperType && !vt.Field(i).IsExported() {
				prevSkip = true
				continue
			} else if w.sqlWrapperType && vt.Field(i).Name == "Valid" {
				prevSkip = true
				continue
			}
			if _, ok := vt.Field(i).Tag.Lookup("protobuf"); ok && !w.protoWrapperType {
				w.protoWrapperType = true
			}
			sv, ok := vt.Field(i).Tag.Lookup("logger")
			switch {
			case ok:
				switch sv {
				case "omit":
					prevSkip = true
					continue
				case "take":
					break
				}
			case !ok && w.opts.Tagged:
				// skip top level untagged
				if w.depth == 1 {
					prevSkip = true
					continue
				}
				if tv, ok := w.takeMap[w.depth]; ok && !tv {
					prevSkip = true
					continue
				}
			}

			if prevSkip {
				prevSkip = false
			}

			if numWritten > 0 {
				_, _ = w.s.Write(commaBytes)
				_, _ = w.s.Write(spaceBytes)
			}

			vt := vt.Field(i)
			if w.s.Flag('+') || w.s.Flag('#') {
				_, _ = w.s.Write([]byte(vt.Name))
				_, _ = w.s.Write(colonBytes)
			}
			w.format(w.unpackValue(v.Field(i)))
			numWritten++
		}
		w.depth--

		if numWritten == 0 && w.depth < 0 {
			_, _ = w.s.Write(filteredBytes)
		}
		_, _ = w.s.Write(closeBraceBytes)
	case reflect.Uintptr:
		getHexPtr(w.s, uintptr(v.Uint()))
	case reflect.UnsafePointer, reflect.Chan, reflect.Func:
		getHexPtr(w.s, v.Pointer())
	// There were not any other types at the time this code was written, but
	// fall back to letting the default fmt package handle it if any get added.
	default:
		format := w.buildDefaultFormat()
		if v.CanInterface() {
			_, _ = fmt.Fprintf(w.s, format, v.Interface())
		} else {
			_, _ = fmt.Fprintf(w.s, format, v.String())
		}
	}
}

func (w *Wrapper) Format(s fmt.State, verb rune) {
	w.s = s

	// Use standard formatting for verbs that are not v.
	if verb != 'v' {
		format := w.constructOrigFormat(verb)
		_, _ = fmt.Fprintf(s, format, w.val)
		return
	}

	if w.val == nil {
		if s.Flag('#') {
			_, _ = s.Write(interfaceBytes)
		}
		_, _ = s.Write(nilAngleBytes)
		return
	}

	if w.opts.Tagged {
		w.buildTakeMap(reflect.ValueOf(w.val), 1)
	}

	w.format(reflect.ValueOf(w.val))
}

// handle special methods like error.Error() or fmt.Stringer interface
func handleMethods(_ *Options, w io.Writer, v reflect.Value) (handled bool) {
	if !v.CanInterface() {
		// not our case
		return false
	}

	if !v.CanAddr() {
		// not our case
		return false
	}

	if v.CanAddr() {
		v = v.Addr()
	}

	// Is it an error or Stringer?
	switch iface := v.Interface().(type) {
	case error:
		defer catchPanic(w, v)
		_, _ = w.Write([]byte(iface.Error()))
		return true
	case fmt.Stringer:
		defer catchPanic(w, v)
		_, _ = w.Write([]byte(iface.String()))
		return true
	}

	return false
}

// getBool outputs a boolean value as true or false to Writer w.
func getBool(w io.Writer, val bool) {
	if val {
		_, _ = w.Write(trueBytes)
	} else {
		_, _ = w.Write(falseBytes)
	}
}

// getInt outputs a signed integer value to Writer w.
func getInt(w io.Writer, val int64, base int) {
	_, _ = w.Write([]byte(strconv.FormatInt(val, base)))
}

// getUint outputs an unsigned integer value to Writer w.
func getUint(w io.Writer, val uint64, base int) {
	_, _ = w.Write([]byte(strconv.FormatUint(val, base)))
}

// getFloat outputs a floating point value using the specified precision,
// which is expected to be 32 or 64bit, to Writer w.
func getFloat(w io.Writer, val float64, precision int) {
	_, _ = w.Write([]byte(strconv.FormatFloat(val, 'g', -1, precision)))
}

// getComplex outputs a complex value using the specified float precision
// for the real and imaginary parts to Writer w.
func getComplex(w io.Writer, c complex128, floatPrecision int) {
	r := real(c)
	_, _ = w.Write(openParenBytes)
	_, _ = w.Write([]byte(strconv.FormatFloat(r, 'g', -1, floatPrecision)))
	i := imag(c)
	if i >= 0 {
		_, _ = w.Write(plusBytes)
	}
	_, _ = w.Write([]byte(strconv.FormatFloat(i, 'g', -1, floatPrecision)))
	_, _ = w.Write(iBytes)
	_, _ = w.Write(closeParenBytes)
}

// getHexPtr outputs a uintptr formatted as hexadecimal with a leading '0x'
// prefix to Writer w.
func getHexPtr(w io.Writer, p uintptr) {
	// Null pointer.
	num := uint64(p)
	if num == 0 {
		_, _ = w.Write(nilAngleBytes)
		return
	}

	// Max uint64 is 16 bytes in hex + 2 bytes for '0x' prefix
	buf := make([]byte, 18)

	// It's simpler to construct the hex string right to left.
	base := uint64(16)
	i := len(buf) - 1
	for num >= base {
		buf[i] = hexDigits[num%base]
		num /= base
		i--
	}
	buf[i] = hexDigits[num]

	// Add '0x' prefix.
	i--
	buf[i] = 'x'
	i--
	buf[i] = '0'

	// Strip unused leading bytes.
	buf = buf[i:]
	_, _ = w.Write(buf)
}

func catchPanic(w io.Writer, _ reflect.Value) {
	if err := recover(); err != nil {
		_, _ = w.Write(panicBytes)
		_, _ = fmt.Fprintf(w, "%v", err)
		_, _ = w.Write(closeParenBytes)
	}
}

func (w *Wrapper) buildDefaultFormat() (format string) {
	buf := bytes.NewBuffer(percentBytes)

	for _, flag := range sf {
		if w.s.Flag(int(flag)) {
			_, _ = buf.WriteRune(flag)
		}
	}

	_, _ = buf.WriteRune('v')

	format = buf.String()
	return format
}

func (w *Wrapper) constructOrigFormat(verb rune) string {
	buf := bytes.NewBuffer(percentBytes)

	for _, flag := range sf {
		if w.s.Flag(int(flag)) {
			_, _ = buf.WriteRune(flag)
		}
	}

	if width, ok := w.s.Width(); ok {
		_, _ = buf.WriteString(strconv.Itoa(width))
	}

	if precision, ok := w.s.Precision(); ok {
		_, _ = buf.Write(precisionBytes)
		_, _ = buf.WriteString(strconv.Itoa(precision))
	}

	_, _ = buf.WriteRune(verb)

	return buf.String()
}

func (w *Wrapper) buildTakeMap(v reflect.Value, depth int) {
	if !v.IsValid() || v.IsZero() {
		return
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			w.buildTakeMap(v.Index(i), depth+1)
		}
		w.takeMap[depth] = true
		return
	case reflect.Struct:
		break
	case reflect.Ptr:
		v = v.Elem()
		if v.Kind() != reflect.Struct {
			w.takeMap[depth] = true
			return
		}
	default:
		w.takeMap[depth] = true
		return
	}

	vt := v.Type()

	for i := 0; i < v.NumField(); i++ {
		sv, ok := vt.Field(i).Tag.Lookup("logger")
		if ok && sv == "take" {
			w.takeMap[depth] = false
		}
		if v.Kind() == reflect.Struct ||
			(v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct) {
			w.buildTakeMap(v.Field(i), depth+1)
		}
	}

	if _, ok := w.takeMap[depth]; !ok {
		w.takeMap[depth] = true
	}
}
