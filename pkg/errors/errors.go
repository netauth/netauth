package errors

import "errors"

var (
	E_NO_ENTITY              = errors.New("No entity matched the given constraints")
	E_ENTITY_UNQUALIFIED     = errors.New("The specified entity does not have sufficient capabilitites")
	E_NO_GROUP               = errors.New("No such group.")
	E_NO_SUCH_DATABASE       = errors.New("No such database!")
	E_NO_SUCH_CRYPTO         = errors.New("No such crypto engine!")
	E_CRYPTO_FAULT           = errors.New("Unrecoverable crypto fault!  Check log for more details.")
	E_CRYPTO_BADAUTH         = errors.New("Bad authentication information!")
	E_BAD_REQUEST            = errors.New("This request is malformed")
)
