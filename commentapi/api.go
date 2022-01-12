package commentapi

import "github.com/lbryio/lbry.go/extras/api"

// Validator for api arguments that should have some checks applied for every request.
type Validator interface {
	Validate() api.StatusError
}
