package utils

import "strconv"

type StringInt string

func (s StringInt) Int() int {
	i, _ := strconv.Atoi(string(s))

	return i
}
