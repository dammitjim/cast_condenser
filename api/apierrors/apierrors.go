package apierrors

import "encoding/json"

// APIError is a standardised error that implements the native
// error interface. Allows more detailed errors to be passed back
// up the chain to the transport layer.
type APIError struct {
	InternalErrorCode int    `json:"condenser_search_error_code"`
	Message           string `json:"message"`
	Details           string `json:"details"`
	HTTPStatusCode    int    `json:"http_status_code"`
}

// implement error interface.
func (e APIError) Error() string {
	jsonBuf, _ := json.Marshal(e)
	return string(jsonBuf)
}

// WithDetails returns a new APIError object with the details added.
func (e APIError) WithDetails(details string) *APIError {
	ae := &e
	ae.Details = details
	return ae
}

var (
	// Generic represents a generic error
	Generic = APIError{InternalErrorCode: 1000, Message: "Something went wrong.", HTTPStatusCode: 500}

	// GenericValidation represents a failure in validation.
	GenericValidation = APIError{InternalErrorCode: 2000, Message: "Validation failed.", HTTPStatusCode: 400}

	// GenericExternal represents a generic error from an external provider.
	GenericExternal = APIError{InternalErrorCode: 3000, Message: "Something went wrong with an external provider.", HTTPStatusCode: 500}

	// ITunesNon200 for when something goes wrong with the itunes request.
	ITunesNon200 = APIError{InternalErrorCode: 3001, Message: "Itunes returned non 200.", HTTPStatusCode: 500}
)
