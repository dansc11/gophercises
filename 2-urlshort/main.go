package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	handler, err := setupHandlers(pathsToUrls, "routes.yml")

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func setupHandlers(routeMap map[string]string, routesYamlPath string) (http.HandlerFunc, error) {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	mapHandler := MapHandler(routeMap, mux)

	// Load the yaml content
	yaml, err := os.ReadFile(routesYamlPath)
	if err != nil {
		return nil, err
	}

	yamlHandler, err := YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		return nil, err
	}

	return yamlHandler, nil
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/hello", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Default route")
}
