package server

import (
	"fmt"
	"net/http"
)

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	//resp := make(map[string]string)
	//resp["message"] = "Hello World"
	//
	//jsonResp, err := json.Marshal(resp)
	//if err != nil {
	//	log.Fatalf("error handling JSON marshal. Err: %v", err)
	//}
	//w.Header().Set("Content-Type", "application/json")
	//_, _ = w.Write(jsonResp)
	fmt.Fprintf(
		w, `
		  ##         .
	## ## ##        ==
 ## ## ## ## ##    ===
/"""""""""""""""""\___/ ===
{                       /  ===-
\______ O           __/
 \    \         __/
  \____\_______/

	
Hello from Docker!

`)
}
