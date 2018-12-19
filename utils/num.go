package utils

import "strconv"

// Atoi func
func Atoi(id string) int {
	i, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	return i
}
