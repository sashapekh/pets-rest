package helper

func GetString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
