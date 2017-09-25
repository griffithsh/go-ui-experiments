package client

import (
	"fmt"
	"net/http"
	"strings"
)

// Resolver is anything that can dereference a map of file paths to their
// contents.
type Resolver interface {
	Resolve(path string) ([]byte, error)
}

// Start the client ui server, which serves everything in the *client* directory.
func Start(host string, resolver Resolver, config map[string]string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Got request %v\n", r.URL)

		if r.URL.Path == "/configuration-values.js" {
			// Inject configuration values
			rows := []string{}
			for k, v := range config {
				rows = append(rows, fmt.Sprintf("'%s':'%v'", k, v))
			}
			template := `var cfg = {%s};`
			w.Write([]byte(fmt.Sprintf(template, strings.Join(rows, ","))))
			return
		}
		b, err := resolver.Resolve(r.URL.Path)
		if err != nil {
			w.Write([]byte(err.Error()))
			fmt.Printf("Resolve %s: %s\n", r.URL.Path, err)
		}
		w.Write(b)
	})

	// Start server
	go http.ListenAndServe(host, nil)
}
