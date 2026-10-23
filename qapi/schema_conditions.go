package qapi

// ParseConditionExpr parses an 'if' condition expression.
func ParseConditionExpr(raw interface{}) *ConditionExpr {
	if raw == nil {
		return nil
	}
	if s, ok := raw.(string); ok {
		return &ConditionExpr{Simple: s}
	}
	if m, ok := raw.(map[string]interface{}); ok {
		cond := &ConditionExpr{}
		if all, ok := m["all"].([]interface{}); ok {
			for _, item := range all {
				if str, ok := item.(string); ok {
					cond.All = append(cond.All, str)
				}
			}
		}
		if any, ok := m["any"].([]interface{}); ok {
			for _, item := range any {
				if str, ok := item.(string); ok {
					cond.Any = append(cond.Any, str)
				}
			}
		}
		if notRaw, ok := m["not"]; ok {
			cond.Not = ParseConditionExpr(notRaw)
		}
		return cond
	}
	return nil
}
