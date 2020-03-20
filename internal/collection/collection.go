package collection

// Tests if any of the items in vs can satisfy function f
// Credit goes to https://gobyexample.com/collection-functions
func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}
