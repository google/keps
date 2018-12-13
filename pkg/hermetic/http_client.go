package hermetic

import (
	"net/http"
)

// packageHttpClient is the HTTP client to use when issuing GitHub requests
var packageHttpClient *http.Client = &http.Client{}
