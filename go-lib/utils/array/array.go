package array

func ContainAny(arr []string, targets []string) bool { // To check if the array contains any of the targets
	targetSet := make(map[string]struct{}, len(targets))
	for _, t := range targets {
		targetSet[t] = struct{}{}
	}

	for _, a := range arr {
		if _, found := targetSet[a]; found {
			return true
		}
	}
	return false
}
