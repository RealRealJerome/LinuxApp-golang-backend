package course

import (
	"strconv"
	"strings"
)

func IntArr2Str(input []int) string {
	if len(input) == 0 {
		return ""
	}
	str := strconv.Itoa(input[0])
	for i := 1; i < len(input); i++ {
		str += "," + strconv.Itoa(input[i])
	}
	return str
}
func Str2IntArr(input string) []int {
	if input == "" {
		return nil
	}
	strSlice := strings.Split(input, ",")
	res := make([]int, 0, len(strSlice))
	for _, item := range strSlice {
		itemInt, _ := strconv.Atoi(item)
		res = append(res, itemInt)
	}
	return res
}
