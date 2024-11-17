package zerogen

import "strings"

// ToSnakeCase converts a string from camelCase or PascalCase to snake_case using a state machine approach
func ToSnakeCase(str string) string {
	var result []rune
	prevLower := false // 是否前一个字符是小写字母
	prevUpper := false // 是否前一个字符是大写字母
	prevDigit := false // 是否前一个字符是数字

	for i, r := range str {
		// 当前字符是大写字母
		if isUpper(r) {
			// 如果前一个是小写字母或数字，说明这里需要插入下划线
			if prevLower || prevDigit {
				result = append(result, '_')
			} else if prevUpper && i > 0 && i < len(str)-1 && isLower(rune(str[i+1])) {
				// 如果前一个是大写字母，当前也是大写，且下一个是小写字母，需要插入下划线
				result = append(result, '_')
			}
			result = append(result, toLower(r))
			prevLower = false
			prevUpper = true
			prevDigit = false
		} else if isDigit(r) { // 当前字符是数字
			// 如果前一个是字母，说明需要插入下划线
			if prevLower || prevUpper {
				result = append(result, '_')
			}
			result = append(result, r)
			prevLower = false
			prevUpper = false
			prevDigit = true
		} else { // 当前字符是小写字母
			result = append(result, r)
			prevLower = true
			prevUpper = false
			prevDigit = false
		}
	}

	return string(result)
}

// Helper function to check if a rune is uppercase
func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

// Helper function to check if a rune is lowercase
func isLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

// Helper function to check if a rune is a digit
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// Helper function to convert an uppercase rune to lowercase
func toLower(r rune) rune {
	if isUpper(r) {
		return r + ('a' - 'A')
	}
	return r
}

// ToLowerCamelCase converts snake_case to lowerCamelCase
func ToLowerCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if i > 0 {
			parts[i] = strings.Title(parts[i])
		} else {
			parts[i] = strings.ToLower(parts[i])
		}
	}
	return strings.Join(parts, "")
}
