package database

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrInvalidDSNAddr      = errors.New("invalid dsn addr")
	ErrInvalidDSNUnescaped = errors.New("dsn must be escaped")
	ErrInvalidDSNNoSlash   = errors.New("dsn must contains slash")
)

type Config struct {
	TLSConfig *tls.Config
	Username  string
	Password  string
	Scheme    string
	Host      string
	Port      string
	Database  string
	Params    []string
}

func (cfg *Config) FormatDSN() string {
	var s strings.Builder

	if len(cfg.Scheme) > 0 {
		s.WriteString(cfg.Scheme + "://")
	}
	// [username[:password]@]
	if len(cfg.Username) > 0 {
		s.WriteString(cfg.Username)
		if len(cfg.Password) > 0 {
			s.WriteByte(':')
			s.WriteString(url.PathEscape(cfg.Password))
		}
		s.WriteByte('@')
	}

	// [host:port]
	if len(cfg.Host) > 0 {
		s.WriteString(cfg.Host)
		if len(cfg.Port) > 0 {
			s.WriteByte(':')
			s.WriteString(cfg.Port)
		}
	}

	// /dbname
	s.WriteByte('/')
	s.WriteString(url.PathEscape(cfg.Database))

	for i := 0; i < len(cfg.Params); i += 2 {
		if i == 0 {
			s.WriteString("?")
		} else {
			s.WriteString("&")
		}
		s.WriteString(cfg.Params[i])
		s.WriteString("=")
		s.WriteString(cfg.Params[i+1])
	}

	return s.String()
}

func ParseDSN(dsn string) (*Config, error) {
	cfg := &Config{}

	// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	// Find last '/' that goes before dbname
	foundSlash := false
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' {
			foundSlash = true
			var j, k int

			// left part is empty if i <= 0
			if i > 0 {
				// Find the first ':' in dsn
				for j = i; j >= 0; j-- {
					if dsn[j] == ':' {
						cfg.Scheme = dsn[0:j]
					}
				}

				// [username[:password]@][host]
				// Find the last '@' in dsn[:i]
				for j = i; j >= 0; j-- {
					if dsn[j] == '@' {
						// username[:password]
						// Find the second ':' in dsn[:j]
						for k = 0; k < j; k++ {
							if dsn[k] == ':' {
								if cfg.Scheme == dsn[:k] {
									continue
								}
								var err error
								cfg.Password, err = url.PathUnescape(dsn[k+1 : j])
								if err != nil {
									return nil, err
								}
								break
							}
						}
						cfg.Username = dsn[len(cfg.Scheme)+3 : k]
						break
					}
				}

				for k = j + 1; k < i; k++ {
					if dsn[k] == ':' {
						cfg.Host = dsn[j+1 : k]
						cfg.Port = dsn[k+1 : i]
						break
					}
				}

			}

			// dbname[?param1=value1&...&paramN=valueN]
			// Find the first '?' in dsn[i+1:]
			for j = i + 1; j < len(dsn); j++ {
				if dsn[j] == '?' {
					parts := strings.Split(dsn[j+1:], "&")
					cfg.Params = make([]string, 0, len(parts)*2)
					for _, p := range parts {
						k, v, found := strings.Cut(p, "=")
						if !found {
							continue
						}
						cfg.Params = append(cfg.Params, k, v)
					}

					break
				}
			}
			var err error
			dbname := dsn[i+1 : j]
			if cfg.Database, err = url.PathUnescape(dbname); err != nil {
				return nil, fmt.Errorf("invalid dbname %q: %w", dbname, err)
			}

			break
		}
	}

	if !foundSlash && len(dsn) > 0 {
		return nil, ErrInvalidDSNNoSlash
	}

	return cfg, nil
}
