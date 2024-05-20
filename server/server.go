package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Run() error {
	r := gin.Default()
	//********************************
	userCtrl.InitRoutes(r)
	productCtrl.InitRoutes(r)
	cartCtrl.InitRoutes(r)
	//********************************
	return r.Run(s.port)
}
