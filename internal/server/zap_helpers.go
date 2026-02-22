package server

import "go.uber.org/zap"

func zapFields(fields []interface{}) []zap.Field {
	var zf []zap.Field

	for i := 0; i < len(fields); i += 2 {
		key := fields[i].(string)
		value := fields[i+1]

		switch v := value.(type) {
		case string:
			zf = append(zf, zap.String(key, v))
		case int:
			zf = append(zf, zap.Int(key, v))
		case int64:
			zf = append(zf, zap.Int64(key, v))
		default:
			zf = append(zf, zap.Any(key, v))
		}
	}

	return zf
}
