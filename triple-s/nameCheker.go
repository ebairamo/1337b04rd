package main

import (
	"regexp"
	"unicode"
)

func checkBucketName(BucketName string) (bool, *ErrorResponse) {
	if BucketName == "" {
		return false, &ErrorResponse{Code: 400, Message: "Name is empty"}
	}
	if len(BucketName) < 3 || len(BucketName) > 63 {
		return false, &ErrorResponse{Code: 400, Message: "Bucket name must be between 3 and 63 characters"}
	}

	for i := 0; i < len(BucketName); i++ {
		char := BucketName[i]
		// Если символ не является буквой, цифрой, дефисом или точкой
		if !(unicode.IsLower(rune(char)) || unicode.IsDigit(rune(char)) || char == '-' || char == '.') {
			return false, &ErrorResponse{Code: 400, Message: "Bucket name can only contain lowercase letters, digits, hyphens (-), and dots (.)"}
		}
	}
	for i := 0; i < len(BucketName)-1; i++ {
		if BucketName[i] == '-' && BucketName[i+1] == '-' {
			return false, &ErrorResponse{Code: 400, Message: "Bucket name cannot contain consecutive hyphens ('--')"}
		}
	}
	re := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if re.MatchString(BucketName) {
		return false, &ErrorResponse{Code: 400, Message: "Bucket name cannot be in IP address format"}
	}
	if BucketName[0] == '-' || BucketName[len(BucketName)-1] == '-' {
		return false, &ErrorResponse{Code: 400, Message: "Bucket name cannot start or end with a hyphen"}
	}

	return true, nil
}
