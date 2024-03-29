package server

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
)

type ApiServer struct {
	svc Service
}

func NewApiServer(svc Service) *ApiServer {
	return &ApiServer{
		svc: svc,
	}
}

func (s *ApiServer) Start(listenAddr string) error {
	http.HandleFunc("/parse", CORS(s.handlePostParseFile))
	return http.ListenAndServe(listenAddr, nil)
}

func (s *ApiServer) handlePostParseFile(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		writeJSON(w, http.StatusNotFound, "")
	case "POST":
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			writeJSON(w, http.StatusBadRequest, err)
			return
		}
		parsedFile, err := s.svc.PostParseFile(body, context.Background())
		if err != nil {
			writeJSON(w, http.StatusBadRequest, err)
		} else {
			writeJSON(w, http.StatusOK, parsedFile)
		}

	}
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")

	switch vv := v.(type) {
	case []byte:
		_, err := w.Write(vv)
		return err
	case error:
		_, err := w.Write([]byte(vv.Error()))
		return err
	default:
		return errors.New("could not process request")
	}
}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}
