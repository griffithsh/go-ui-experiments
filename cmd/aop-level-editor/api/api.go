package api

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
)

// state is a database and image files referenced by that database.
type state struct {
	dir string // dir that we've opened that contains a database file and image files
	db  database
}

// newState constructs a state
func newState(db *sql.DB) *state {
	return &state{}
}

// entry is the entry point for all *api* queries. There's something similar in
// the client directory.
func (s *state) entry(w http.ResponseWriter, r *http.Request) {
	switch r.Method + r.URL.Path {
	case "POST/inject":
		// On POSTS to /inject, the request body is interpreted as SQL, and
		// executed on the DB.
		if b, err := ioutil.ReadAll(r.Body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("read request body: %s", err)))
		} else {
			result, err := s.db.execute(string(b))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("db execute: %s", err)))
			}
			w.Write([]byte(fmt.Sprintf("%s", result)))
		}
	case "POST/open":
		// On POSTs to /open, a new database/resource folder is opened to
		// perform operations on.

		// TODO ...
	default:
		// Otherwise, return a client error.
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Printf("how exciting, we got a %s for %s!\n", r.Method, r.URL.Path)
	}
}

// Start the api server, and return the url to access it with.
func Start(port int, db *sql.DB) string {
	host := fmt.Sprintf("127.0.0.1:%d", port)

	// Start server
	s := newState(db)
	go http.ListenAndServe(host, http.HandlerFunc(s.entry))
	return "http://" + host
}
