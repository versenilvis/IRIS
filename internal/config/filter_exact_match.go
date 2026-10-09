package config

import "fmt"

type FilterExactMatchMode int

const (
	FilterExactMatchAuto FilterExactMatchMode = iota
	FilterExactMatchOn
	FilterExactMatchOff
)

func (f *FilterExactMatchMode) UnmarshalTOML(val any) error {
	switch v := val.(type) {
	case bool:
		if v {
			*f = FilterExactMatchOn
		} else {
			*f = FilterExactMatchOff
		}
		return nil
	case string:
		if v == "auto" {
			*f = FilterExactMatchAuto
			return nil
		}
	}
	return fmt.Errorf("filter-exact-match must be a boolean or \"auto\"")
}

func (f FilterExactMatchMode) MarshalTOML() ([]byte, error) {
	switch f {
	case FilterExactMatchAuto:
		return []byte(`"auto"`), nil
	case FilterExactMatchOn:
		return []byte("true"), nil
	case FilterExactMatchOff:
		return []byte("false"), nil
	default:
		return nil, fmt.Errorf("filter-exact-match: invalid value %d", f)
	}
}
