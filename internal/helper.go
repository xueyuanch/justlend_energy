package internal

// Contains is a helper function convenient to inspect if
// an element inside a variable arguments.
func Contains[T comparable](e T, es ...T) bool {
	for _, v := range es {
		if v == e {
			return true
		}
	}
	return false
}

// IsEmpty returns a boolean indicating whether the input string is empty or not
func IsEmpty(r string) bool { return len(r) == 0 }

func IsValidAddress(addr string) bool {
	if len(addr) != 34 {
		return false
	}
	if string(addr[0:1]) != "T" {
		return false
	}
	return DecodeCheck(addr) != nil
}
