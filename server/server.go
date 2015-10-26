package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()

	addApiTasks(r)

	// mocks...
	r.HandleFunc("/user/{id}", wrapJsonHandler(GetUserHandler)).Methods("GET", "OPTIONS")

	http.Handle("/", r)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, "mocks/user.json")
	return nil
}
