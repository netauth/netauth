package db

import (
	"errors"
)

var (
	// ErrUnknownEntity is returned for requests to load an entity
	// that does not exist.
	ErrUnknownEntity = errors.New("The specified entity does not exist")

	// ErrUnknownGroup is returned for requests to load a group
	// that does not exist.
	ErrUnknownGroup = errors.New("The specified group does not exist")

	// ErrUnknownDatabase is returned for an attempt to create a
	// new database that hasn't been registered.
	ErrUnknownDatabase = errors.New("The specified database does not exist")

	// ErrInternalError is used for all other errors that occur
	// within a database implementation.
	ErrInternalError = errors.New("The database has encountered an internal error")
)
