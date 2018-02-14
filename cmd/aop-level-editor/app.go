package main

import (
	"database/sql"
	"flag"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/griffithsh/go-ui-experiments/cmd/aop-level-editor/api"
	"github.com/griffithsh/go-ui-experiments/cmd/aop-level-editor/client"
	"github.com/griffithsh/go-ui-experiments/cmd/aop-level-editor/client/static"
	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

var (
	clientHost = "127.0.0.1:4000"
)

func newStaticResolver() *staticResolver {
	return &staticResolver{}

}

type staticResolver struct {
}

func (r *staticResolver) Resolve(path string) ([]byte, error) {
	return static.Resolve(path)
}

func main() {
	// Parse flags
	flag.Parse()

	// Create astilectron
	var err error
	var a *astilectron.Astilectron
	astilog.SetLogger(astilog.New(astilog.FlagConfig()))
	if a, err = astilectron.New(astilectron.Options{
		AppName:            "Astilectron",
		AppIconDefaultPath: "",
		AppIconDarwinPath:  "",
		BaseDirectoryPath:  "./",
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating new astilectron failed"))
	}
	defer a.Close()
	a.HandleSignals()

	// Start astilectron
	if err = a.Start(); err != nil {
		astilog.Fatal(errors.Wrap(err, "starting failed"))
	}

	// Start client code server
	db, err := sql.Open("sqlite3", "testfile")
	defer db.Close()
	if err != nil {
		astilog.Fatal(errors.Wrap(err, "opening db failed"))
		return
	}
	apiPath := api.Start(8181, db)
	client.Start(clientHost, newStaticResolver(), map[string]string{"api-url": apiPath})

	// Create window
	var w *astilectron.Window
	if w, err = a.NewWindow("http://"+clientHost, &astilectron.WindowOptions{
		Title:  astilectron.PtrStr("aop-level-editor"),
		Center: astilectron.PtrBool(true),
		Height: astilectron.PtrInt(480),
		Width:  astilectron.PtrInt(640),
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "new window failed"))
	}
	if err = w.Create(); err != nil {
		astilog.Fatal(errors.Wrap(err, "creating window failed"))
	}

	// Blocking pattern
	a.Wait()
}
