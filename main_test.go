package main
import ("testing";"net/http/httptest";"strings";"encoding/json")
func TestShorten(t *testing.T){
 h:=newHandler();w:=httptest.NewRecorder();h.ServeHTTP(w,httptest.NewRequest("POST","/shorten",strings.NewReader(`{"url":"https://example.com"}`)))
 if w.Code!=201{t.Fatal(w.Code)}
 var result map[string]string;json.Unmarshal(w.Body.Bytes(),&result)
 w=httptest.NewRecorder();h.ServeHTTP(w,httptest.NewRequest("GET",result["path"],nil))
 if w.Code!=302||w.Header().Get("Location")!="https://example.com"{t.Fatal(w.Code,w.Header())}
}
func TestInvalid(t *testing.T){w:=httptest.NewRecorder();newHandler().ServeHTTP(w,httptest.NewRequest("POST","/shorten",strings.NewReader(`{"url":"javascript:alert(1)"}`)));if w.Code!=400{t.Fatal(w.Code)}}
