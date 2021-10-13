package main

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

type RedirectHandler struct {
	pathMap  map[string]string
	fallback http.Handler
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path

	for path, redirectTo := range h.pathMap {
		if path == uri {
			http.Redirect(w, r, redirectTo, http.StatusTemporaryRedirect)
		}
	}

	h.fallback.ServeHTTP(w, r)
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	redirectHandler := RedirectHandler{
		pathMap:  pathsToUrls,
		fallback: fallback,
	}

	return redirectHandler.ServeHTTP
}

type redirectDetails struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var redirects []redirectDetails

	if err := yaml.Unmarshal(yml, &redirects); err != nil {
		return nil, err
	}

	var yamlPathMap = make(map[string]string)

	for _, redirect := range redirects {
		yamlPathMap[redirect.Path] = redirect.Url
	}

	return MapHandler(yamlPathMap, fallback), nil
}
