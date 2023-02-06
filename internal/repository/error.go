package error_response

import (
	"embed"
	"encoding/json"

	"github.com/labstack/echo/v4"
)

//go:embed error_properties.json
var error_properties embed.FS
var errorData map[string]ErrorData

type ResponseError struct {
	Message    string         `json:"message"`
	Path       []string       `json:"path,omitempty"`
	Extensions map[string]any `json:"extensions,omitempty"`

	Redacted bool  `json:"-"`
	Err      error `json:"-"`
}

type ErrorData struct {
	Body   any `json:"body,omitempty"`
	Status int `json:"status,omitempty"`
}

func constructErrorData() {
	filename := "error_properties.json"
	error_properties, _ := error_properties.ReadFile(filename)
	json.Unmarshal(error_properties, &errorData)
}

func CreateError(c echo.Context, error_key string) error {
	if errorData == nil {
		constructErrorData()
	}

	return c.JSON(errorData[error_key].Status, errorData[error_key].Body)
}
