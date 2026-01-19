package helper

import "encoding/json"

func NormalizeToIntMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		result[key] = normalizeValueToInt(value)
	}

	return result
}

func normalizeValueToInt(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		// Recursively normalize nested maps
		return NormalizeToIntMap(v)

	case []interface{}:
		// Normalize array elements
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = normalizeValueToInt(item)
		}
		return result

	case json.Number:
		// Convert json.Number to int
		if i, err := v.Int64(); err == nil {
			return int(i)
		}
		// If not int, try float
		if f, err := v.Float64(); err == nil {
			return int(f)
		}
		return v

	case float64:
		// Convert float64 to int
		return int(v)

	case float32:
		return int(v)

	case int64:
		return int(v)

	case int32:
		return int(v)

	case uint:
		return int(v)

	case uint64:
		return int(v)

	case uint32:
		return int(v)

	case int:
		// Already int
		return v

	case string:
		// Keep string as-is
		return v

	case bool:
		// Keep bool as-is
		return v

	default:
		// Unknown type, return as-is
		return v
	}
}
