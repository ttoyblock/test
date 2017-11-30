package type_conv

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrNil = errors.New("Type Converter: nil returned.")

func ToInt(origin interface{}) (int, error) {
	switch v := origin.(type) {
	case int64:
		x := int(v)
		if int64(x) != v {
			return 0, strconv.ErrRange
		}
		return x, nil
	case []byte:
		x, err := strconv.Atoi(string(v))
		return x, err
	case nil:
		return 0, ErrNil
	}
	return 0, fmt.Errorf("Type Converter: unexpected type for int, got type %T", origin)
}

func ToInt64(origin interface{}) (int64, error) {
	switch v := origin.(type) {
	case int64:
		return v, nil
	case []byte:
		x, err := strconv.ParseInt(string(v), 10, 64)
		return x, err
	case nil:
		return 0, ErrNil
	}
	return 0, fmt.Errorf("Type Converter: unexpected type for int64, got type %T", origin)
}

var errNegativeInt = errors.New("ype Converter: unexpected value for Uint64")

func ToUint64(origin interface{}) (uint64, error) {
	switch v := origin.(type) {
	case int64:
		if v < 0 {
			return 0, errNegativeInt
		}
	case []byte:
		x, err := strconv.ParseUint(string(v), 10, 64)
		return x, err
	case nil:
		return 0, ErrNil
	}
	return 0, fmt.Errorf("Type Converter: unexpected type for Uint64, got type %T", origin)
}

func ToFloat64(origin interface{}) (float64, error) {
	switch v := origin.(type) {
	case []byte:
		x, err := strconv.ParseFloat(string(v), 64)
		return x, err
	case nil:
		return 0, ErrNil
	case string:
		r, err := strconv.ParseFloat(string(v), 64)
		return r, err
	}
	return 0, fmt.Errorf("Type Converter: unexpected type for Float64, got type %T", origin)
}

func ToString(origin interface{}) (string, error) {
	switch v := origin.(type) {
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case []byte:
		return string(v), nil
	case string:
		return v, nil
	case nil:
		return "", ErrNil
	}
	return "", fmt.Errorf("Type Converter: unexpected type for String, got type %T", origin)
}

func ToByteSlice(origin interface{}) ([]byte, error) {
	switch v := origin.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case nil:
		return nil, ErrNil
	}
	return nil, fmt.Errorf("Type Converter: unexpected type for []byte, got type %T", origin)
}

func ToBool(origin interface{}) (bool, error) {
	switch v := origin.(type) {
	case int:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case []byte:
		return strconv.ParseBool(string(v))
	case nil:
		return false, ErrNil
	}
	return false, fmt.Errorf("Type Converter: unexpected type for Bool, got type %T", origin)
}
