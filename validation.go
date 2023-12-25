package main

func SliceUnique(s []string) bool {
	vals := map[string]struct{}{}
	for _, v := range s {
		_, exists := vals[v]
		if exists {
			return false
		}
		vals[v] = struct{}{}
	}
	return true
}
