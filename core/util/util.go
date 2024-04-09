package util

import "strings"

func PascalCase(snake string) string {
	// 将snake_case字符串按下划线分割成多个部分
	parts := strings.Split(snake, "_")
	// 对于所有部分，将其转换为大写
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}
	// 将所有部分拼接起来
	return strings.Join(parts, "")
}

func SQLTypeToStructType(sqlType string) string {
	switch sqlType {
	case "bigint":
		return "int64"
	case "int", "tinyint":
		return "int32"
	case "text", "varchar":
		return "string"
	case "decimal", "double", "float":
		return "float64"
	case "datetime", "timestamp":
		return "jsontime.JsonTime"
	default:
		return ""
	}
}
