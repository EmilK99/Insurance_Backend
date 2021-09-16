package api

import (
	"errors"
	"net/http"
)

var ErrLowCancellationRate = errors.New("Low cancellation rate")

// ErrorResponse https://developer.paypal.com/docs/api/errors/
type ErrorResponse struct {
	Response        *http.Response        `json:"-"`
	Name            string                `json:"name"`
	DebugID         string                `json:"debug_id"`
	Message         string                `json:"message"`
	InformationLink string                `json:"information_link"`
	Details         []ErrorResponseDetail `json:"details"`
}

// ErrorResponseDetail struct
type ErrorResponseDetail struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
	Links []Link `json:"link"`
}

// Link struct
type Link struct {
	Href        string `json:"href"`
	Rel         string `json:"rel,omitempty"`
	Method      string `json:"method,omitempty"`
	Description string `json:"description,omitempty"`
	Enctype     string `json:"enctype,omitempty"`
}
