package shared

type Variable[T string | int64 | float64 | bool] string

type Config interface {
	GetString(key Variable[string]) string
	GetInt64(key Variable[int64]) int64
	GetFloat64(key Variable[float64]) float64
	GetBool(key Variable[bool]) bool
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
