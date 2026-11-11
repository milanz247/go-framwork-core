package support

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type M = map[string]any

func Str(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func Integer(v any) (int64, error) {
	switch n := v.(type) {
	case int64:
		return n, nil
	case int:
		return int64(n), nil
	case float64:
		// FormatFloat preserves whole-number digits without exponent notation.
		return strconv.ParseInt(strconv.FormatFloat(n, 'f', -1, 64), 10, 64)
	case json.Number:
		return n.Int64()
	case string:
		return strconv.ParseInt(n, 10, 64)
	default:
		return 0, fmt.Errorf("expected integer, received %T", v)
	}
}

func Number(v any) int64 { n, _ := Integer(v); return n }
