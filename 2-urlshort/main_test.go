package main

import (
	"io/fs"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSetupHandlers(t *testing.T) {
	// GIVEN
	yamlPath := "test-routes.yml"
	yamlContent := `- path: /yaml-a
  url: /yaml-b
- path: /yaml-c
  url: /yaml-d
`

	if err := os.WriteFile(yamlPath, []byte(yamlContent), fs.ModePerm); err != nil {
		log.Fatal(err)
	}

	defer os.Remove(yamlPath)

	routesMap := map[string]string{
		"/map-a": "/map-b",
		"/map-c": "/map-d",
	}

	// WHEN
	handler, err := setupHandlers(routesMap, yamlPath)

	// THEN
	if err != nil {
		t.Fatal("Unexpected error from setupHandlers")
	}

	testHandlerRedirect(t, handler, "/map-a", "/map-b")
	testHandlerRedirect(t, handler, "/map-c", "/map-d")
	testHandlerRedirect(t, handler, "/yaml-a", "/yaml-b")
	testHandlerRedirect(t, handler, "/yaml-c", "/yaml-d")
	testHandlerOk(t, handler, "/hello")
	testHandlerOk(t, handler, "/")
}

func testHandlerRedirect(t *testing.T, handler http.HandlerFunc, route string, expectRedirect string) {
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	testHandler := http.HandlerFunc(handler)

	testHandler.ServeHTTP(rr, req)

	wantStatus := http.StatusTemporaryRedirect
	if gotStatus := rr.Code; wantStatus != gotStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", gotStatus, wantStatus)
	}

	wantLocation := expectRedirect
	gotLocation := rr.HeaderMap.Get("Location")
	if gotLocation != wantLocation {
		t.Errorf("handler returned unexpected redirect location: got %s want %s",
			gotLocation, wantLocation)
	}
}

func testHandlerOk(t *testing.T, handler http.HandlerFunc, route string) {
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	testHandler := http.HandlerFunc(handler)

	testHandler.ServeHTTP(rr, req)

	wantStatus := http.StatusOK
	if gotStatus := rr.Code; wantStatus != gotStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", gotStatus, wantStatus)
	}
}
