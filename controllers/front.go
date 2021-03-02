package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

// Registers userController and assigns to available endpoints
func RegisterControllers() {
	uc := newUserController()

	http.Handle("/users", *uc)
	http.Handle("/users/", *uc)
}

// Encodes data into json and writes using supplied Writer interface
func encodeResponseAsJSON(data interface{}, writer io.Writer)  {
	encoder := json.NewEncoder(writer)
	encoder.Encode(data)
}

// Writes status code into header and feedback into body using supplied Writer interface.
func writeResponse(writer http.ResponseWriter, statusCode int, feedback string) {
	writer.WriteHeader(statusCode)
	writer.Write([]byte(feedback))
}