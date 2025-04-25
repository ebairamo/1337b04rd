package main

import (
	"encoding/xml"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/xml")
	localBucket := ErrorResponse{
		Code:    code,
		Message: message,
	}
	x, _ := xml.MarshalIndent(localBucket, "", " ")
	w.Write(x)
}
