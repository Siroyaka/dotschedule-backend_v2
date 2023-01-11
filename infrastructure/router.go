package infrastructure

import (
	"log"
	"net/http"
)

type Router struct {
	port string
}

func NewRouter(port string) Router {
	return Router{port: port}
}

func (r *Router) Run() {
	log.Printf("Port:%s", r.port)
	log.Fatal(http.ListenAndServe(r.port, nil))
}

func (_ *Router) SetHandle(route string, handler http.Handler) {
	http.Handle(route, handler)
}

func (_ *Router) SetHandleFunc(route string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(route, handler)
}
