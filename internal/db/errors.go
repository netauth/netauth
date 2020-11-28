package db

import (
	"errors"
)

var (
	// ErrUnknownEntity is returned for requests to load an entity
	// that does not exist.
	ErrUnknownEntity = errors.New("the specified entity does not exist")

	// ErrUnknownGroup is returned for requests to load a group
	// that does not exist.
	ErrUnknownGroup = errors.New("the specified group does not exist")

	// ErrUnknownDatabase is returned for an attempt to create a
	// new database that hasn't been registered.
	ErrUnknownDatabase = errors.New("the specified database does not exist")

	// ErrInternalError is used for all other errors that occur
	// within a database implementation.
	ErrInternalError = errors.New("the database has encountered an internal error")

	// ErrBadSearch is returned when a search request cannot be
	// filled for some reason.
	ErrBadSearch = errors.New("the provided SearchRequest is invalid")

	// ErrNoValue is returned when no value exists for a given key.
	ErrNoValue = errors.New("no value exists")
)
