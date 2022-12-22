package server

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/jiangklijna/web-shell/lib"
)

// Version WebShell Server current version
const Version = "1.0"

// Server Response header[Server]
const Server = "web-shell-" + Version

// WebShellServer Main Server
type WebShellServer struct {
	http.ServeMux
	cache lib.ExpiredMap
}

// StaticHandler reserved for static_gen.go
var StaticHandler http.Handler

// Init WebShell. register handlers
func (s *WebShellServer) Init(ContentPath string, Command ...string) {
	if StaticHandler == nil {
		StaticHandler = HTMLDirHandler()
	}
	s.cache = *lib.NewExpiredMap()

	s.Handle(ContentPath+"/", s.upgrade(ContentPath, StaticHandler))
	s.Handle(ContentPath+"/cmd/", s.upgrade(ContentPath, VerifyHandler(s.PathVerifyFunc, ConnectionHandler(Command...))))
	s.Handle(ContentPath+"/login", s.upgrade(ContentPath, LoginHandler(s.PasswordVerifyFunc, s.PathHandleFunc)))
}

// packaging and upgrading http.Handler
func (s *WebShellServer) upgrade(ContentPath string, h http.Handler) http.Handler {
	return LoggingHandler(GetMethodHandler(ContentPathHandler(ContentPath, h)))
}

// Run WebShell server
func (s *WebShellServer) Run(https bool, port, crt, key, rootcrt string) {
	var err error
	server := &http.Server{Addr: ":" + port, Handler: s}
	if https {
		if rootcrt != "" {
			server.TLSConfig = &tls.Config{
				ClientCAs:  lib.ReadCertPool(rootcrt),
				ClientAuth: tls.RequireAndVerifyClientCert,
			}
		}
		err = server.ListenAndServeTLS(crt, key)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (s *WebShellServer) PasswordVerifyFunc(username, password string) bool {
	//TODO check from database
	return true
}
func (s *WebShellServer) PathVerifyFunc(path string) bool {
	ok, _ := s.cache.Get(path)
	return ok
}
func (s *WebShellServer) PathHandleFunc(path string) {
	s.cache.Set(path, 1, 10)
}
