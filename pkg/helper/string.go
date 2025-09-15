package helper

/*
GetString safely converts an interface{} to a string.
If the conversion is not possible, it returns an empty string.
*/
func GetString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
