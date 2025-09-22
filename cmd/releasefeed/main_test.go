package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)
	return rr
}

func TestHandleProduct(t *testing.T) {
	s := CreateNewServer()
	s.MountHandlers()

	req, _ := http.NewRequest("GET", "/feeds/mediawiki", nil)
	response := executeRequest(req, s)
	require.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	require.Contains(t, body, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	require.Contains(t, body, "<title>mediawiki</title>")
	require.Contains(t, body, "<id>tag:,2021-09-30:mediawiki:1.31.16</id>")
}

func TestHhandleCycle(t *testing.T) {
	s := CreateNewServer()
	s.MountHandlers()

	req, _ := http.NewRequest("GET", "/feeds/mediawiki/1.31", nil)
	response := executeRequest(req, s)
	require.Equal(t, http.StatusOK, response.Code)

	body := response.Body.String()
	require.Contains(t, body, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	require.Contains(t, body, "<title>mediawiki (1.31)</title>")
	require.Contains(t, body, "<id>tag:,2021-09-30:mediawiki:1.31.16</id>")
}
