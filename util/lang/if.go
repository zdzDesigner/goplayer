package lang

// If3 简单三元表达
func If3(ok bool, a, b interface{}) interface{} {
	if ok {
		return a
	}
	return b
}
