package trans

func StringNvl(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Int64Nvl(s *int64) int64 {
	if s == nil {
		return 0
	}
	return *s
}

func BoolToString(b bool) string {
	if b {
		return "ON"
	}
	return "OFF"
}
