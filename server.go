package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

type Server struct {
	leds *ws
}

func NewServerInstance() *Server {
	lights := NewLEDLights()

	server := &Server{
		leds: lights,
	}

	InintTwitchWS(server)

	server.run()

	return server
}

func (s *Server) run() {
	log.SetFlags(0)

	http.HandleFunc("/", s.serveTemplate)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/favicon.ico", http.StripPrefix("/static/", fs))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/api/values", s.handleValues)

	// Start the server on localhost port 8082 and log any errors
	log.Println("http server started on :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (s *Server) serveTemplate(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)
	if path == "\\" {
		path = "\\index.html"
	}

	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", path)

	log.Println(path)

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Println("Error parsing files!")
		return
	}
	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println("Error executing the template!")
	}
}

func (s *Server) handleValues(w http.ResponseWriter, r *http.Request) {
	log.Println("Values!")
	switch r.Method {
	case "", "GET":
		if err := json.NewEncoder(w).Encode(LEDValues{
			Display:    s.leds.display,
			Brightness: s.leds.brightness,
			Colors: []string{
				"FF000",
			},
		}); err != nil {
			log.Printf("error: %v\n", err)
		}
	case "POST":
		bytes, _ := io.ReadAll(r.Body)

		var values LEDValues
		if err := json.Unmarshal(bytes, &values); err != nil {
			log.Printf("error: %v\n", err)
		}

		s.leds.display = values.Display
		s.leds.brightness = values.Brightness
	}
}
