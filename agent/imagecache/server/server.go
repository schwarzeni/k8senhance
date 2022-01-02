package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/schwarzeni/k8senhance/agent/imagecache/job"
	"github.com/schwarzeni/k8senhance/config"
)

type Server struct {
	conf *config.Config
	r    *mux.Router
}

func (s *Server) Run() error {
	go job.NewHealthJob(s.conf).Run()
	go job.NewMetricJob(s.conf).Run()
	srv := &http.Server{
		Handler: s.r,
		Addr:    s.conf.Agent.Imagecache.Addr,
	}
	return srv.ListenAndServe()
}

func NewServer(config *config.Config) *Server {
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(">>>>>>> [debug]", r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})
	server := &Server{r: r, conf: config}
	HandleService(server)
	HandleProxy(server)
	return server
}
