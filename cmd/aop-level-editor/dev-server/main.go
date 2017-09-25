package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/griffithsh/go-ui-experiments/cmd/aop-level-editor/api"
	"github.com/griffithsh/go-ui-experiments/cmd/aop-level-editor/client"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbFlag      string
	webRootFlag string
)

func init() {
	flag.StringVar(&dbFlag, "db", "", "path of the database file to edit")
	flag.StringVar(&webRootFlag, "root", "", "path to the root of the client website to serve")
}

func newFileResolver(path string) *fileResolver {
	return &fileResolver{clientPath: path}
}

type fileResolver struct {
	clientPath string
}

func (fr *fileResolver) Resolve(path string) ([]byte, error) {
	if path == "/" {
		// Default '/' to index.html
		return ioutil.ReadFile(fr.clientPath + "/index.html")
	}
	// Otherwise treat other paths as a request to the file system
	return ioutil.ReadFile(fr.clientPath + "/" + path)
}

func main() {
	flag.Parse()

	// Start api server, then start client server.
	db, err := sql.Open("sqlite3", dbFlag)
	defer db.Close()
	if err != nil {
		log.Fatalf("could not open db: %s\n", err)
	}
	apiPath := api.Start(8001, db)
	client.Start("127.0.0.1:8002", newFileResolver(webRootFlag), map[string]string{"api-url": apiPath})

	fmt.Printf("Listening on :%d\n", 8002)
	// TODO: Try to pop a browser window.

	// Block until we get an interrupt.
	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		wg.Done()
	}()
	wg.Wait()
}
