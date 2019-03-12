package main

import (
	"flag"
	tmpl "html/template"
	"log"
	"net/http"
	"strconv"

	playground "github.com/keksboter/setlx2python_playground"
)

func main() {
	mode := flag.String("mode", "prod", "run mode, dev or prod")
	port := flag.Int("port", 80, "port which the webserver listens on")
	setlxRunnerURL := flag.String("setlxRunnerURL", "https://setlx.dotcookie.me/run", "the url to the setlx code runner")
	flag.Parse()

	log.Printf("Starting setlx2python playground server in %s mode on port %d\n", *mode, *port)

	// load page html template
	template, err := tmpl.ParseFiles("www/index.html")
	if err != nil {
		log.Fatalln(err)
		return
	}

	handler := playground.NewRequestHandler(template, *setlxRunnerURL)

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(*port),
		Handler: playground.CreateRouter(handler),
	}

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Println("webserver failed:", err)
	}
}
