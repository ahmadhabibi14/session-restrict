package converter

import "strconv"

func Float64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}
