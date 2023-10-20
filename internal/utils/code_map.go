package utils

var CodeMap map[int64]string

func init() {
	ResponseStatusCode()
}

func ResponseStatusCode() {
	CodeMap = map[int64]string{
		200: "Successfully",
		1001: "Wrong with the data format trasnfer from POST",
		1002: "Failed to execute SQL syntax",
		1003: "No search this user",
		1004: "Invalid password",
		1005: "Failed to create session value JSON data",
	}
}
