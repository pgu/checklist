package server

import (
	"appengine/datastore"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type badRequest struct{ error }

func wrapJsonHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// - setup response
		w.Header().Set("Content-Type", "application/json")

		//
		// EXEC handler
		err := f(w, r)

		//
		// POST OPS
		// - handle an error
		if err == nil {
			return
		}

		if err == datastore.ErrNoSuchEntity {
			http.Error(w, "entity not found", http.StatusNotFound)
			return
		}

		switch err.(type) {
		case badRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "oops", http.StatusInternalServerError)
		}
	}
}

func parseID(r *http.Request) (int64, error) {
	txt, ok := mux.Vars(r)["id"]
	if !ok {
		return 0, fmt.Errorf("task id not found")
	}
	return strconv.ParseInt(txt, 10, 0)
}
