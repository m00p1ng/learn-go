package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/m00p1ng/learn-go/nomadcoin/blockchain"
)

const (
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.GetBlockChain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockChain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(aPort int) {
	handler := http.NewServeMux()
	port := fmt.Sprintf("localhost:%d", aPort)
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.go.html"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.go.html"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("Listening on %s", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
