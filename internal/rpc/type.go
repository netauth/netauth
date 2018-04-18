package rpc

type EntityTree interface {
	ValidateSecret(string, string) error
	MakeBootstrap(string, string)
	DisableBootstrap()
}

type NetAuthServer struct {
	Tree EntityTree
}
