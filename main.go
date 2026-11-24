package main
import ("crypto/rand";"encoding/hex";"encoding/json";"net/http";"net/url";"sync";"log")
type store struct { mu sync.RWMutex; links map[string]string }
func newHandler() http.Handler {
 s:= &store{links:make(map[string]string)}
 m:=http.NewServeMux()
 m.HandleFunc("POST /shorten",func(w http.ResponseWriter,r *http.Request){
  var input struct{ URL string }
  if json.NewDecoder(http.MaxBytesReader(w,r.Body,4096)).Decode(&input)!=nil {http.Error(w,"Invalid JSON",400);return}
  u,e:=url.Parse(input.URL);if e!=nil||u.Hostname()==""||(u.Scheme!="http"&&u.Scheme!="https"){http.Error(w,"Invalid URL",400);return}
  s.mu.Lock();defer s.mu.Unlock()
  if len(s.links)>=10000 {http.Error(w,"Capacity reached",503);return}
  var code string
  for { b:=make([]byte,6);if _,e:=rand.Read(b);e!=nil{http.Error(w,"Unavailable",500);return};code=hex.EncodeToString(b);if _,ok:=s.links[code];!ok{break} }
  s.links[code]=u.String()
  w.Header().Set("Content-Type","application/json");w.WriteHeader(201);json.NewEncoder(w).Encode(map[string]string{"path":"/"+code})
 })
 m.HandleFunc("GET /{code}",func(w http.ResponseWriter,r *http.Request){s.mu.RLock();target,ok:=s.links[r.PathValue("code")];s.mu.RUnlock();if !ok{http.NotFound(w,r);return};http.Redirect(w,r,target,302)})
 return m
}
func main(){log.Fatal(http.ListenAndServe("127.0.0.1:8080",newHandler()))}
