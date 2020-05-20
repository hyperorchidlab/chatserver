package api

import "net/http"

type UserCommand struct {

}


func (uc *UserCommand)ServeHTTP(w http.ResponseWriter, r *http.Request)  {

}

