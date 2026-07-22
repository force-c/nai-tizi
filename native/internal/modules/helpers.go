package modules

func containsCode(codes []string, expected string) bool {
	if len(codes) == 0 {
		return true
	}
	for _, code := range codes {
		if code == expected {
			return true
		}
	}
	return false
}
