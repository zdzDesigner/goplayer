package lang

// MapAssign ..
func MapAssign(m1, m2 map[string]string) map[string]string {
	newmp := map[string]string{}
	for k, v := range m1 {
		newmp[k] = v
	}
	for k, v := range m2 {
		newmp[k] = v
	}
	return newmp
}
