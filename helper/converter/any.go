package converter

import (
	"session-restrict/src/lib/logger"
	"strconv"
	"time"

	"github.com/goccy/go-json"
)

func AnyToInt64(x any) int64 {
	if x == nil {
		return 0
	}
	if val, ok := x.(int64); ok {
		return val
	}
	switch v := x.(type) {
	case int:
		return int64(v)
	case uint:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case time.Duration:
		return int64(v)
	case *int:
		if v != nil {
			return int64(*v)
		}
	case *uint:
		if v != nil {
			return int64(*v)
		}
	case *int8:
		if v != nil {
			return int64(*v)
		}
	case *int16:
		if v != nil {
			return int64(*v)
		}
	case *int32:
		if v != nil {
			return int64(*v)
		}
	case *int64:
		if v != nil {
			return *v
		}
	case *uint8:
		if v != nil {
			return int64(*v)
		}
	case *uint16:
		if v != nil {
			return int64(*v)
		}
	case *uint32:
		if v != nil {
			return int64(*v)
		}
	case *uint64:
		if v != nil {
			return int64(*v)
		}
	case *float32:
		if v != nil {
			return int64(*v)
		}
	case *float64:
		if v != nil {
			return int64(*v)
		}
	case bool:
		if v {
			return 1
		}
		return 0
	case []byte:
		if val, err := strconv.ParseInt(string(v), 10, 64); err == nil {
			return val
		}
		if val, err := strconv.ParseFloat(string(v), 64); err == nil {
			return int64(val)
		}
	case string:
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return val
		}
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return int64(val)
		}
	case *any:
		if v != nil {
			return AnyToInt64(*v)
		}
	default:
		return 0
	}
	return 0
}

func AnyToUint64(x any) uint64 {
	if x == nil {
		return 0
	}
	if val, ok := x.(uint64); ok {
		return val
	}
	switch v := x.(type) {
	case int:
		return uint64(v)
	case uint:
		return uint64(v)
	case int8:
		return uint64(v)
	case int16:
		return uint64(v)
	case int32:
		return uint64(v)
	case int64:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint32:
		return uint64(v)
	case uint64:
		return v
	case float32:
		return uint64(v)
	case float64:
		return uint64(v)
	case time.Duration:
		return uint64(v)
	case *int:
		if v != nil {
			return uint64(*v)
		}
	case *uint:
		if v != nil {
			return uint64(*v)
		}
	case *int8:
		if v != nil {
			return uint64(*v)
		}
	case *int16:
		if v != nil {
			return uint64(*v)
		}
	case *int32:
		if v != nil {
			return uint64(*v)
		}
	case *int64:
		if v != nil {
			return uint64(*v)
		}
	case *uint8:
		if v != nil {
			return uint64(*v)
		}
	case *uint16:
		if v != nil {
			return uint64(*v)
		}
	case *uint32:
		if v != nil {
			return uint64(*v)
		}
	case *uint64:
		if v != nil {
			return *v
		}
	case *float32:
		if v != nil {
			return uint64(*v)
		}
	case *float64:
		if v != nil {
			return uint64(*v)
		}
	case bool:
		if v {
			return 1
		}
		return 0
	case []byte:
		if val, err := strconv.ParseInt(string(v), 10, 64); err == nil {
			return uint64(val)
		}
		if val, err := strconv.ParseFloat(string(v), 64); err == nil {
			return uint64(val)
		}
	case string:
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return uint64(val)
		}
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return uint64(val)
		}
	case *any:
		if v != nil {
			return AnyToUint64(*v)
		}
	default:
		return 0
	}
	return 0
}

func AnyToMap(x any) map[string]any {
	if x == nil {
		return map[string]any{}
	}
	val, ok := x.(map[string]any)
	if ok {
		return val
	}
	val, ok = x.(map[string]any)
	if !ok {
		return map[string]any{}
	}

	return val
}

func AnyToString(x any) string {
	if x == nil {
		return ""
	}
	val, ok := x.(string)
	if !ok {
		return ""
	}

	return val
}

func AnyToJsonPretty(any any) string {
	res, err := json.MarshalIndent(any, ``, `  `)
	if err != nil {
		logger.Log.Error(err)
		return ""
	}
	return string(res)
}
