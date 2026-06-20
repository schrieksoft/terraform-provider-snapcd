package integration

// Shared across the integration data source + the supply/event resources.

var integrationEndpoint = "/Integration"

// lookup finds a key case-insensitively (tolerating camelCase or PascalCase JSON).
func lookup(m map[string]interface{}, key string) (interface{}, bool) {
	if v, ok := m[key]; ok {
		return v, true
	}
	for k, v := range m {
		if len(k) == len(key) && toLowerFirst(k) == toLowerFirst(key) {
			return v, true
		}
	}
	return nil, false
}

func toLowerFirst(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	if b[0] >= 'A' && b[0] <= 'Z' {
		b[0] += 'a' - 'A'
	}
	return string(b)
}

func getStr(m map[string]interface{}, key string) string {
	if v, ok := lookup(m, key); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := lookup(m, key); ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := lookup(m, key); ok {
		if mm, ok := v.(map[string]interface{}); ok {
			return mm
		}
	}
	return nil
}
