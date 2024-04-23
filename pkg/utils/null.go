package utils

func NullString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
