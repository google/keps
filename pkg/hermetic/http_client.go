package hermetic

import (
	"net/http"
)

var packageHttpClient *http.Client = &http.Client{}
