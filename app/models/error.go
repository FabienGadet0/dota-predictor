package models

// Error handling
type Error struct {
	Code    int
	Message string
	Data    interface{}
}
