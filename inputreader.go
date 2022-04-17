package prombench

import "net/http"

// Converts input data to *http.Request.
type InputReader interface {
	Read() (*http.Request, error)
}
