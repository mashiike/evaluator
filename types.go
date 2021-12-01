package evaluator

import "strconv"

func isBothStrings(v1, v2 interface{}) (s1, s2 string, ok bool) {
	s1, ok = isString(v1)
	if !ok {
		return
	}
	s2, ok = isString(v2)
	return
}

func isString(v interface{}) (string, bool) {
	switch v := v.(type) {
	case string:
		return v, true
	case rune:
		return string(v), true
	default:
		return "", false
	}
}

func isBool(v interface{}) (b, ok bool) {
	b, ok = v.(bool)
	return
}

func isBothBools(v1, v2 interface{}) (b1, b2, ok bool) {
	b1, ok = isBool(v1)
	if !ok {
		return
	}
	b2, ok = isBool(v2)
	return
}

func isBothRealNumbers(v1, v2 interface{}) (n1, n2 float64, ok bool) {
	n1, ok = isRealNumber(v1)
	if !ok {
		return
	}
	n2, ok = isRealNumber(v2)
	return
}

func isRealNumber(v interface{}) (float64, bool) {
	switch v := v.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

func asNumber(v interface{}) (float64, bool) {
	if n, ok := isRealNumber(v); ok {
		return n, true
	}
	if s, ok := isString(v); ok {
		f, err := strconv.ParseFloat(s, 64)
		return f, err == nil
	}
	return 0.0, false
}

func asString(v interface{}) (string, bool) {
	if s, ok := isString(v); ok {
		return s, true
	}
	if n, ok := isRealNumber(v); ok {
		return strconv.FormatFloat(n, 'f', -1, 64), true
	}

	return "", false
}

func asBool(v interface{}) (bool, bool) {
	if b, ok := isBool(v); ok {
		return b, ok
	}
	if n, ok := isRealNumber(v); ok {
		return n == 0, true
	}
	if s, ok := isString(v); ok {
		b, err := strconv.ParseBool(s)
		return b, err == nil
	}
	return false, false
}
