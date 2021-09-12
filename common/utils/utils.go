package utils

// String returns pointer to s.
func String(s string) *string {
	return &s
}

// Int returns a pointer to i.
func Int(i int) *int {
	return &i
}

// StringSlice converts a slice of string to a
// slice of *string
func StringSlice(elements ...string) []*string {
	var res []*string
	for _, element := range elements {
		e := element
		res = append(res, &e)
	}
	return res
}
