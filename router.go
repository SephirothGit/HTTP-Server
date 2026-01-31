package server

import "net/http"

func NewRouter(updateStatus http.HandleFunc) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/orders/", updateStatus)

	return mux
}
