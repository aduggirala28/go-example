package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Config struct {
	addr string
}


func home(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err!=nil {
		log.Printf("Could not render home page html %s",err.Error() )
		http.Error(w, "Internal Server Error", 500)
	}

	err = ts.Execute(w, nil)
	if err!=nil{
		log.Println("Error occured while executing template")
		http.Error(w, "Internal Server Error", 500)
	}
	w.Write([]byte("Welcome to SnippetBox Homepage"))
}

func showSnippet(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err!=nil || id<1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "This is the snippet chosen %d ", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(405)
		w.Write([]byte("HTTP Method not supported for this route."))
		return
	}
	w.Write([]byte("Snippet created."))
}

func main(){

	conf := new(Config)
	flag.StringVar(&conf.addr,"addr", ":4003", "HTTP Address of the Server")
	flag.Parse()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fmt.Println("Listening on port", conf.addr)
	err := http.ListenAndServe(conf.addr, mux)
	log.Fatal(err)

}
