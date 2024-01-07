package apps

import "net/http"

type HttpServer struct {
	server *http.Server
}

func NewHttpServer(addr string, handler http.Handler) *HttpServer {
	return &HttpServer{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *HttpServer) Run() error {
	return s.server.ListenAndServe()
}

func (s *HttpServer) RunTLS(certFile, keyFile string) error {
	return s.server.ListenAndServeTLS(certFile, keyFile)
}
