package util

func StringAddr(s string) *string {
	if s == "" {
		return nil
	} else {
		return &s
	}
}

func StringValue(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}

func DefaultString(s string, defaultValue string) string {
	if s == "" {
		return defaultValue
	}
	return s
}

func NotNullString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func ConditionString(condition bool, f1 interface{}, f2 interface{}) string {
	if condition {
		if f, ok := f1.(func() string); ok {
			return f()
		} else {
			return f1.(string)
		}
	} else {
		if f, ok := f2.(func() string); ok {
			return f()
		} else {
			return f2.(string)
		}
	}
}

func SubString(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		return ""
	}

	end := begin + length
	if end > lth {
		end = lth
	}

	return string(rs[begin:end])
}
