package setlx2python_playground

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/oxtoacart/bpool"
)

// RequestHandler answers all http requests
type RequestHandler struct {
	indexTemplate  *template.Template
	setlxRunnerURL string
	bufpool        *bpool.BufferPool
}

// NewRequestHandler creates a new handler with a own buffer pool
func NewRequestHandler(indexTemplate *template.Template, setlxRunnerURL string) *RequestHandler {
	return &RequestHandler{
		indexTemplate:  indexTemplate,
		setlxRunnerURL: setlxRunnerURL,
		bufpool:        bpool.NewBufferPool(64),
	}
}

// PageData is metadata required by the template
type PageData struct {
	SetlXCode  string
	PythonCode string
	URL        string
	Host       string
}

func (h *RequestHandler) index(w http.ResponseWriter, r *http.Request) {
	code := "print(\"Hello setlX\");"

	source := r.URL.Query().Get("source")
	if len(source) > 0 {
		content, err := fetchURL(source)
		if err == nil {
			code = *content
		}
	}
	output, err := transpile([]byte(code))
	if err != nil {
		log.Println(err)
		http.Error(w, "Can not transpile code", http.StatusInternalServerError)
		return
	}
	err = h.indexTemplate.Execute(w, PageData{
		SetlXCode:  code,
		PythonCode: output.Code,
		Host:       r.URL.Host,
	})
	if err != nil {
		log.Println("error executing template: ", err)
	}
}

func (h *RequestHandler) transpile(w http.ResponseWriter, r *http.Request) {
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
}

func (h *RequestHandler) runSetlX(w http.ResponseWriter, r *http.Request) {
	res, err := http.Post(h.setlxRunnerURL, "text", r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	w.Write(body)
}

func (h *RequestHandler) runPython(w http.ResponseWriter, r *http.Request) {
	code, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Can not read code", http.StatusInternalServerError)
		return
	}
	output, err := executePython(code)
	if err != nil {
		log.Println(err)
		http.Error(w, "Can not execute code", http.StatusInternalServerError)
		return
	}
	w.Write(output)
}
