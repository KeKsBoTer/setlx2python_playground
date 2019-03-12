package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	tmpl "html/template"

	"github.com/gorilla/mux"
)

// PageData is metadata required by the template
type PageData struct {
	SetlXCode  string
	PythonCode string
	URL        string
	Host       string
}

func main() {
	mode := flag.String("mode", "prod", "run mode, dev or prod")
	port := flag.Int("port", 80, "port which the webserver listens on")
	flag.Parse()

	log.Printf("Starting setlx2python playground server in %s mode on port %d\n", *mode, *port)

	// load page html template
	template, err := tmpl.ParseFiles("www/index.html")
	if err != nil {
		log.Fatalln(err)
		return
	}

	router := mux.NewRouter()
	router.StrictSlash(true)

	router.Path("/").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if *mode == "dev" {
			template, err = tmpl.ParseFiles("www/index.html")
			if err != nil {
				log.Fatalln(err)
				return
			}
		}
		source := r.URL.Query().Get("source")
		if len(source) > 0 {
			log.Println(source)
		}
		code := "print(\"Hello setlX\");"
		output, err := transpile([]byte(code))
		if err != nil {
			log.Println(err)
			http.Error(w, "Can not transpile code", http.StatusInternalServerError)
			return
		}
		err = template.Execute(w, PageData{
			SetlXCode:  code,
			PythonCode: output.Code,
			Host:       r.URL.Host,
		})
		if err != nil {
			log.Println("error executing template: ", err)
		}
	})

	router.Path("/transpile").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Can not read code", http.StatusInternalServerError)
			return
		}
		response, err := transpile(code)
		if err != nil {
			log.Println(err)
			http.Error(w, "Can not transpile code", http.StatusInternalServerError)
			return
		}
		output, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			http.Error(w, "Can not create response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
	})
	// serve static files
	fileHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("www/static")))
	router.PathPrefix("/static/").Handler(fileHandler)

	server := &http.Server{Addr: ":" + strconv.Itoa(*port), Handler: router}

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Println("webserver failed:", err)
	}
}
