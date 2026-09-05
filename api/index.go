package handler

import (
	"net/http"

	"github.com/MishraShardendu22/pkg/server"
)

var appServer = server.NewAppServer()
var httpHandler = appServer.Handler()

func Handler(w http.ResponseWriter, r *http.Request) {
	httpHandler.ServeHTTP(w, r)
}
