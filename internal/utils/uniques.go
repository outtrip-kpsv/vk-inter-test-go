package utils

func uniqueInts(slice []int) []int {
	encountered := make(map[int]bool)
	unique := []int{}

	for _, v := range slice {
		if !encountered[v] {
			encountered[v] = true
			unique = append(unique, v)
		}
	}

	return unique
}

func UniqueValues(m map[int][]int) []int {
	values := []int{}

	for _, v := range m {
		values = append(values, v...)
	}

	return uniqueInts(values)
}
