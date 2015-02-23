package helper

import "errors"

func IsUint(data interface{}) bool {
	switch data.(type) {
	case uint64, uint32, uint16, uint8, uint:
		return true
	}
	return false
}

func ToUint64(data interface{}) (uint64, error) {
	switch data.(type) {
	case uint64:
		return data.(uint64), nil
	case uint32:
		return uint64(data.(int32)), nil
	case uint16:
		return uint64(data.(int16)), nil
	case uint8:
		return uint64(data.(int8)), nil
	case uint:
		return uint64(data.(int)), nil
	}
	return 0, errors.New("Invalid conversion to uint64")
}

// Helper method, converting all int and float types in an interface to a float64.
func ToFloat64(data interface{}) (float64, error) {
	switch data.(type) {
	case float64:
		return data.(float64), nil
	case float32:
		return float64(data.(float32)), nil
	case int64:
		return float64(data.(int64)), nil
	case int32:
		return float64(data.(int32)), nil
	case int16:
		return float64(data.(int16)), nil
	case int8:
		return float64(data.(int8)), nil
	case int:
		return float64(data.(int)), nil
	}
	return 0, errors.New("Invalid conversion to float64")
}

func ToString(data interface{}) (string, error) {
	switch data.(type) {
	case string:
		return data.(string), nil
	case []byte:
		return string(data.([]byte)), nil
	case []rune:
		return string(data.([]rune)), nil
	}

	return "", errors.New("Invalid conversion to string")
}
