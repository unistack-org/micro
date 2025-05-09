package time

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/goccy/go-yaml"
)

type Duration int64

func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, fmt.Errorf(`time: invalid duration "%s"`, s)
	}

	var p int
	var hours int
loop:
	for i, r := range s {
		switch r {
		case 's', 'm':
			break loop
		case 'h':
			d, err := strconv.Atoi(s[p:i])
			if err != nil {
				return 0, fmt.Errorf(`time: invalid duration "%s"`, s)
			}
			hours += d
			p = i + 1
		case 'd':
			d, err := strconv.Atoi(s[p:i])
			if err != nil {
				return 0, fmt.Errorf(`time: invalid duration "%s"`, s)
			}
			hours += d * 24
			p = i + 1
		case 'y':
			n, err := strconv.Atoi(s[p:i])
			if err != nil {
				return 0, fmt.Errorf(`time: invalid duration "%s"`, s)
			}
			var d int
			for j := n - 1; j >= 0; j-- {
				d += time.Date(time.Now().Year()+j, time.December, 31, 0, 0, 0, 0, time.Local).YearDay()
			}
			hours += d * 24
			p = i + 1
		}
	}

	return time.ParseDuration(fmt.Sprintf("%dh%s", hours, s[p:]))
}

func (d Duration) MarshalYAML() (interface{}, error) {
	return time.Duration(d).String(), nil
}

func (d *Duration) UnmarshalYAML(data []byte) error {
	var v interface{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		dv, err := ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(dv)
		return nil
	default:
		return fmt.Errorf("invalid duration")
	}
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
		dv, err := ParseDuration(value)
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
