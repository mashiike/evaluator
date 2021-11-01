package evaluator

func isBothStrings(v1, v2 interface{}) (s1, s2 string, ok bool) {
	s1, ok = v1.(string)
	if !ok {
		return
	}
	s2, ok = v2.(string)
	return
}

func isBothBools(v1, v2 interface{}) (b1, b2, ok bool) {
	b1, ok = v1.(bool)
	if !ok {
		return
	}
	b2, ok = v2.(bool)
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
