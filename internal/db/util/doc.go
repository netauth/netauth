// Package util handles synthetic functions that aren't implemented
// natively on a datastore.  These functions are almost by definition
// slower than if the datastore implemented an intelligent version
// natively, but they save work and test complexity in the simpler
// datastores.
package util
