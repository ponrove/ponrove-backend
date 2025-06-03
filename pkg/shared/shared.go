package shared

import "errors"

var ErrMissingVariable = errors.New("missing configuration variables")

type Variable[T string | int64 | float64 | bool] string

type Config interface {
	GetString(key Variable[string]) string
	GetInt64(key Variable[int64]) int64
	GetFloat64(key Variable[float64]) float64
	GetBool(key Variable[bool]) bool
	ConfigurationKeysRegistered(keys ...any) error
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

// missingVariableError is an error type that holds a list of missing configuration variable keys.
type missingVariableError struct {
	Keys []string
}

// Error implements the error interface for missingVariableError.
func (e missingVariableError) Error() string {
	return "missing configuration variables: " + formatKeys(e.Keys)
}

// formatKeys formats the keys into a string for error messages. If no keys are provided, it returns "none".
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

// checkKey checks if the provided key exists in the configuration. It uses type assertion to determine the type of the
// key and checks the corresponding map in the configuration struct.
func (c ConfigImpl) checkKey(key any) bool {
	var exists bool
	switch any(key).(type) {
	case Variable[string]:
		_, exists = c.Str[key.(Variable[string])]
	case Variable[int64]:
		_, exists = c.I64[key.(Variable[int64])]
	case Variable[float64]:
		_, exists = c.F64[key.(Variable[float64])]
	case Variable[bool]:
		_, exists = c.Bool[key.(Variable[bool])]
	}

	return exists
}

// ConfigurationKeysRegistered checks if all provided keys are registered in the configuration. To ensure that the
// client of the package have taken all required keys into consideration when building the configuration object.
func (c ConfigImpl) ConfigurationKeysRegistered(keys ...any) error {
	var missingKeys []string
	for _, key := range keys {
		if exists := c.checkKey(key); !exists {
			missingKeys = append(missingKeys, string(key.(Variable[string])))
		}
	}

	if len(missingKeys) > 0 {
		return missingVariableError{Keys: missingKeys}
	}

	return nil
}
