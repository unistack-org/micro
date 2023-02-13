package time

import (
	"encoding/json"
	"fmt"
	"time"
)

type Duration int64

func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, fmt.Errorf(`time: invalid duration "` + s + `"`)
	}

	//var sb strings.Builder
	/*
		for i, r := range s {
			switch r {
			case 'd':
				n, err := strconv.Atoi(s[idx:i])
				if err != nil {
					return 0, errors.New("time: invalid duration " + s)
				}
				s[idx:i] = fmt.Sprintf("%d", n*24)
			default:
				sb.WriteRune(r)
			}
		}
	*/
	var td time.Duration
	var err error
	switch s[len(s)-1] {
	case 's', 'm', 'h':
		td, err = time.ParseDuration(s)
	case 'd':
		if td, err = time.ParseDuration(s[:len(s)-1] + "h"); err == nil {
			td *= 24
		}
	case 'y':
		if td, err = time.ParseDuration(s[:len(s)-1] + "h"); err == nil {
			year := time.Date(time.Now().Year(), time.December, 31, 0, 0, 0, 0, time.Local)
			days := year.YearDay()
			td *= 24 * time.Duration(days)
		}
	}

	return td, err
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		dv, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(dv)
		return nil
	default:
		return fmt.Errorf("invalid duration")
	}
}

/*
func (d Duration) MarshalYAML() (interface{}, error) {
	return nil, nil
}

func (d Duration) UnmarshalYAML(fn func(interface{}) error) error
*/
