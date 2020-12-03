package cart

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gorilla/mux"
// )

// func setup() {
// 	m = mux.NewRouter()
// 	SetRoutes(m)

// 	respRec = httptest.NewRecorder()
// }

// func TestGet400(t *testing.T) {
// 	setup()
// 	//Testing get of non existent question type
// 	req, err = http.NewRequest("GET", "/questions/1/SC", nil)
// 	if err != nil {
// 		t.Fatal("Creating 'GET /questions/1/SC' request failed!")
// 	}

// 	m.ServeHTTP(respRec, req)

// 	if respRec.Code != http.StatusBadRequest {
// 		t.Fatal("Server error: Returned ", respRec.Code, " instead of ", http.StatusBadRequest)
// 	}
// }