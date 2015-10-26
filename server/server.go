package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()

	addApiItems(r)

	http.Handle("/", r)
}
