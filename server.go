package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ps: proxy server
var ps *server

type server struct {
	h  *handler
	r  *httprouter.Router
	sc *scheduler
}

func defaultServer() *server {
	return &server{
		h:  defaultHandler(),
		r:  httprouter.New(),
		sc: nil,
	}
}

func (s *server) run() {
	miniLog.info("Server listening on ", config.PPool.Address)
	s.r.GET("/", s.h.Index)
	s.r.GET("/all", s.h.All)
	s.r.GET("/random", s.h.Random)
	s.r.GET("/size", s.h.Size)

	// start scheduler
	go s.sc.Start()

	// run a server
	miniLog.fatal(http.ListenAndServe(config.PPool.Address, s.r))
	miniLog.info("run after server")
}

func (s *server) receive(spider *Spider) {
	s.sc = getScheduler(spider)
}
