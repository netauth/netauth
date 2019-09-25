package rpc2

// New returns a ready to use server implementation.
func New(r Refs) *Server {
	return &Server{
		Service: r.TokenService,
		Manager: r.Tree,
	}
}
