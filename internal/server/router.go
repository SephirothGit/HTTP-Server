package server

import "net/http"

func NewRouter(updateStatus http.HandlerFunc) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/orders/", updateStatus)

	return mux
}
