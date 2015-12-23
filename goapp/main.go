package main

//go:generate go-bindata -pkg main -o bindata.go assets/
import (
	"log"
	"net/http"

	"github.com/fhermeni/testr/datastore"
	"github.com/fhermeni/testr/rest"
)

//go:generate esc -o static.go -pkg main assets
func init() {
	log.Println("Starting the daemon")

	backend := datastore.NewProvider()
	//st := backend.NewStore(nil)
	//st.Authorize("btrplace", "scheduler", "AAAAAAsKKFKZJKJNZ9298428")

	//Connect the assets
	rest, err := rest.NewEndPoints(backend, "assets")
	if err != nil {
		log.Fatalln(err.Error())
	}

	http.Handle("/assets/", http.FileServer(FS(false)))
	http.Handle("/", rest.Router)
}
