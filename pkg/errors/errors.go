package errors

import "errors"

var (
	E_DUPLICATE_ID           = errors.New("An entity with that ID already exists!")
	E_DUPLICATE_UIDNUMBER    = errors.New("An entity with that number already exists!")
	E_NO_ENTITY              = errors.New("No entity matched the given constraints")
	E_NO_CAPABILITY          = errors.New("The specified capability does not exist!")
	E_ENTITY_UNQUALIFIED     = errors.New("The specified entity does not have sufficient capabilitites")
	E_ENTITY_BADAUTH         = errors.New("The provided entity could not successfully authenticate.")
	E_NO_GROUP               = errors.New("No such group.")
	E_DUPLICATE_GROUP_ID     = errors.New("A group with that ID already exists!")
	E_DUPLICATE_GROUP_NUMBER = errors.New("A group with that number already exists!")
	E_NO_SUCH_DATABASE       = errors.New("No such database!")
	E_NO_SUCH_CRYPTO         = errors.New("No such crypto engine!")
	E_CRYPTO_FAULT           = errors.New("Unrecoverable crypto fault!  Check log for more details.")
	E_CRYPTO_BADAUTH         = errors.New("Bad authentication information!")
	E_EXISTING_EXPANSION     = errors.New("An expansion already exists for that group")
	E_BAD_REQUEST            = errors.New("This request is malformed")
)
