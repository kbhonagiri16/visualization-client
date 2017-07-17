package common

import (
	"encoding/json"
	"net/http"
)

// InvalidOpenstackToken Error
type InvalidOpenstackToken struct{}

// InvalidUsernamePassword Error
type InvalidUsernamePassword struct{}

func (e InvalidOpenstackToken) Error() string {
	return "provided openstack token is not valid"
}

func (e InvalidUsernamePassword) Error() string {
	return "Invalid username or password"
}

// UserDataError means that data, provided by user is not valid
type UserDataError struct {
	Msg string
}

// NewUserDataError return new UserDataError
func NewUserDataError(msg string) UserDataError {
	return UserDataError{msg}
}

func (e UserDataError) Error() string {
	return e.Msg
}

// ClientError means that error occured with client we depend on
type ClientError struct {
	Msg string
}

// NewClientError return new UserDataError
func NewClientError(msg string) ClientError {
	return ClientError{msg}
}

func (e ClientError) Error() string {
	return e.Msg
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

// WriteErrorToResponse writes given params to http response
func WriteErrorToResponse(w http.ResponseWriter, code int, message string,
	details string) {

	w.WriteHeader(code)
	payload, _ := json.Marshal(&errorResponse{code, message, details})
	w.Write(payload)
}
