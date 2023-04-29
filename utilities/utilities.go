package utilities

func GetMaxUint(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}

func GetMinUint(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

func Clamp(value, low, high uint) uint {
	if high < low {
		low, high = high, low
	}
	return GetMinUint(high, GetMaxUint(low, value))
}
