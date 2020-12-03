package routes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)
var m *mux.Router
var req *http.Request
var err error
var rec *httptest.ResponseRecorder

func setup() {
	m = mux.NewRouter()
	SetRoutes(m)

	rec = httptest.NewRecorder()
}
func init() {
	setup() 
}
func Test_Get_Info_Valid(t *testing.T) { 
	req, err = http.NewRequest("GET", "/v1/info", nil)
	if err != nil {
		t.Fatal("Failed to create 'GET /v1/info' request.")
	}
	m.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatal("Server error: Returned ", rec.Code, " expected: ", http.StatusOK)
	}
}

func Test_Login_Valid(t *testing.T) {  
	loginPayload := []byte(`{"username":"nik", "password":"nik"}`) 
	req, err = http.NewRequest("POST", "/v1/login", bytes.NewBuffer(loginPayload))
	if err != nil {
		t.Fatal("Failed to create 'POST /v1/login' request.")
	}
	m.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatal("Server error: Returned ", rec.Code, " expected: ", http.StatusOK)
	}
}

