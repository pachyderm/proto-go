package protoprocess

type apiServer struct {
	processor Processor
	server    Server
}

func newAPIServer(processor Processor, server Server) *apiServer {
	return &apiServer{processor, server}
}

func (a *apiServer) Do(apiDoServer API_DoServer) error {
	return a.server.Handle(a.processor, apiDoServer)
}
