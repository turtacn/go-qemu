package qapi

// ParseFeatures parses features slice from raw features interface.
func ParseFeatures(raw []interface{}) []string {
	var res []string
	for _, f := range raw {
		if s, ok := f.(string); ok {
			res = append(res, s)
		} else if m, ok := f.(map[string]interface{}); ok {
			if name, ok := m["name"].(string); ok {
				res = append(res, name)
			}
		}
	}
	return res
}
