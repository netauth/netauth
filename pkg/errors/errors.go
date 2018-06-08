package errors

import "errors"

var (
	E_ENTITY_UNQUALIFIED     = errors.New("The specified entity does not have sufficient capabilitites")
	E_BAD_REQUEST            = errors.New("This request is malformed")
)
