package encoding

import "strings"

const (
	base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// EncodeBase62 encodes a number to base62 string
func EncodeBase62(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var result []byte
	for num > 0 {
		result = append(result, base62Chars[num%62])
		num /= 62
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// DecodeBase62 decodes a base62 string to a number
func DecodeBase62(str string) int64 {
	var result int64
	for _, char := range str {
		result = result*62 + int64(strings.IndexByte(base62Chars, byte(char)))
	}
	return result
}
