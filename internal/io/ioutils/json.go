package ioutils

import (
  "encoding/json"
  "net/http"
)

func RespJson(w http.ResponseWriter, answer interface{}) {
  js, err := json.MarshalIndent(answer, "", " ")
  if err != nil {
    return
  }
  _, _ = w.Write(js)
}
