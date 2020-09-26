package models

// Response handling
type Response struct {
	Code    int
	Message string
	Data    interface{}
}
