package rpc

type EntityManager interface {
	NewEntity(string, string, string, int32, string) error
	MakeBootstrap(string, string)
	DisableBootstrap()
	DeleteEntity(string, string, string) error
	ChangeSecret(string, string, string, string) error
	ValidateEntitySecretByID(string, string) error
}

type NetAuthServer struct {
	EM EntityManager
}
