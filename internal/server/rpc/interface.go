package rpc

type EntityManager interface {
	NewEntity(string, string, string, int32, string) error
	DeleteEntity(string, string, string) error
	ChangeSecret(string, string, string, string) error
	ValidateSecret(string, string) error
	MakeBootstrap(string, string)
	DisableBootstrap()
}

type NetAuthServer struct {
	EM EntityManager
}
