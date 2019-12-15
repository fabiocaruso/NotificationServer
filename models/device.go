package models

type Device struct {
	ID string			`json:"id"`
	Name string			`json:"name"`
	Os string			`json:"os"`
	Services map[string]interface{}	`json:"services"`
}
