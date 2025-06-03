package shared

import "errors"

var ErrMissingVariable = errors.New("missing configuration variables")

type Variable[T string | int64 | float64 | bool] string

type Config interface {
	GetString(key Variable[string]) string
	GetInt64(key Variable[int64]) int64
	GetFloat64(key Variable[float64]) float64
	GetBool(key Variable[bool]) bool
	ConfigExists(keys ...any) error
}

type ConfigImpl struct {
	Str  map[Variable[string]]string
	I64  map[Variable[int64]]int64
	F64  map[Variable[float64]]float64
	Bool map[Variable[bool]]bool
}

func (c ConfigImpl) GetString(key Variable[string]) string {
	if value, exists := c.Str[key]; exists {
		return value
	}
	return ""
}

func (c ConfigImpl) GetInt64(key Variable[int64]) int64 {
	if value, exists := c.I64[key]; exists {
		return value
	}
	return 0
}

func (c ConfigImpl) GetFloat64(key Variable[float64]) float64 {
	if value, exists := c.F64[key]; exists {
		return value
	}
	return 0.0
}

func (c ConfigImpl) GetBool(key Variable[bool]) bool {
	if value, exists := c.Bool[key]; exists {
		return value
	}
	return false
}

type missingVariableError struct {
	Keys []string
}

func (e missingVariableError) Error() string {
	return "missing configuration variables: " + formatKeys(e.Keys)
}

func formatKeys(keys []string) string {
	if len(keys) == 0 {
		return "none"
	}
	result := ""
	for i, key := range keys {
		if i > 0 {
			result += ", "
		}
		result += string(key)
	}
	return result
}

var _ error = (*missingVariableError)(nil)

func (c ConfigImpl) ConfigExists(keys ...any) error {
	var missingKeys []string
	for _, key := range keys {
		switch any(key).(type) {
		case Variable[string]:
			if _, exists := c.Str[key.(Variable[string])]; !exists {
				missingKeys = append(missingKeys, string(key.(Variable[string])))
				continue
			}
		case Variable[int64]:
			if _, exists := c.I64[key.(Variable[int64])]; !exists {
				missingKeys = append(missingKeys, string(key.(Variable[int64])))
				continue
			}
		case Variable[float64]:
			if _, exists := c.F64[key.(Variable[float64])]; !exists {
				missingKeys = append(missingKeys, string(key.(Variable[float64])))
				continue
			}
		case Variable[bool]:
			if _, exists := c.Bool[key.(Variable[bool])]; !exists {
				missingKeys = append(missingKeys, string(key.(Variable[bool])))
				continue
			}
		}
	}

	if len(missingKeys) > 0 {
		return missingVariableError{Keys: missingKeys}
	}

	return nil
}
