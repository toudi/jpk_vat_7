package common

func LineIsEmpty(line []string) bool {
	for _, item := range line {
		if item != "" {
			return false
		}
	}
	return true
}
