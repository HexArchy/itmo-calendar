package http

import (
	"net/http/pprof"

	"github.com/gorilla/mux"
)

// registerDebugRoutes adds debug and profiling endpoints to the router.
func (s *Server) registerDebugRoutes(router *mux.Router) {
	debug := router.PathPrefix("/debug").Subrouter()
	debug.HandleFunc("/pprof/", pprof.Index)
	debug.HandleFunc("/pprof/profile", pprof.Profile)
	debug.HandleFunc("/pprof/symbol", pprof.Symbol)
	debug.HandleFunc("/pprof/trace", pprof.Trace)
	debug.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	debug.Handle("/pprof/heap", pprof.Handler("heap"))
	debug.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	debug.Handle("/pprof/block", pprof.Handler("block"))
	debug.Handle("/pprof/allocs", pprof.Handler("allocs"))
	debug.Handle("/pprof/mutex", pprof.Handler("mutex"))

	s.logger.Info("Registered debug routes")
}
