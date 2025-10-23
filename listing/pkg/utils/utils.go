package utils

func Join(slice []string, sep string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
