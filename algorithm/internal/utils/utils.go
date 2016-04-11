package utils

func IsSameSignAndNonZero(a, b float64) bool {
	if a == 0 || b == 0 {
		return false
	}
	return (a < 0 && b < 0) || (a > 0 && b > 0)
}
