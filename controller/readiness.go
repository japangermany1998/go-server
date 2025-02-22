package controller

import "net/http"

func HandleReadiness(writer http.ResponseWriter, request *http.Request) {
	request.Header.Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = writer.Write([]byte("OK"))
}
