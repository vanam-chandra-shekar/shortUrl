package server

import (
	"log"
	"net/http"
)

type Srv struct {
	mux    *http.ServeMux
	server *http.Server
	Addr   string
	Port   string
}

func NewServer(addr string, port string) *Srv {
	_mux := http.NewServeMux()
	return &Srv{
		mux: _mux,
		server: &http.Server{
			Addr:    addr + ":" + port,
			Handler: _mux,
		},
		Addr: addr,
		Port: port,
	}
}

func (s *Srv) Register(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *Srv) Use(middleware func(http.Handler) http.Handler) {
	s.server.Handler = middleware(s.server.Handler)
}

func (s *Srv) Run() {
	log.Printf("[%s:%s] server started.", s.Addr, s.Port)
	s.server.ListenAndServe()
}
