package utils

const alphabet = "123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func HashFunction(num int) string {
	if num == 0 {
		return string(alphabet[0])
	}

	base := 61
	result := ""

	for num > 0 {
		rem := num % base
		result = string(alphabet[rem]) + result
		num = num / base
	}

	return result
}
