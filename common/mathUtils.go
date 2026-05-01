package common

func Abs(v int) int {
	if v < 0 {
		return -v
	}

	return v
}

func Sign(n float32) float32 {
	if n > 0 {
		return 1
	}

	if n < 0 {
		return -1
	}

	return 0
}
