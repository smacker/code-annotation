package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type renderFunc func(r *http.Request) response

type response struct {
	statusCode int
	Data       interface{} `json:"data,omitempty"`
	Errors     []errObj    `json:"errors,omitempty"`
}

type errObj struct {
	Title string `json:"title"`
	// we don't have internal code of error or description for now
}

func render(fn renderFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, fn(r))
	}
}

func writeResponse(w http.ResponseWriter, resp response) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(resp.statusCode)
	if resp.statusCode == http.StatusNoContent {
		return
	}
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// Helpers to making response

func respOK(d interface{}) response {
	return response{
		statusCode: http.StatusOK,
		Data:       d,
	}
}

func respErr(statusCode int, err error) response {
	return respStringErr(statusCode, err.Error())
}

func respStringErr(statusCode int, msg string) response {
	// we shouldn't expose real error client
	if statusCode >= 500 {
		logrus.Error(msg)
		msg = http.StatusText(statusCode)
	}
	return response{
		statusCode: statusCode,
		Errors:     []errObj{errObj{Title: msg}},
	}
}

// because respErr(http.StatusInternalServerError, err) is too long
func resp500(err error) response {
	return respErr(http.StatusInternalServerError, err)
}
func resp500f(format string, a ...interface{}) response {
	return respErr(http.StatusInternalServerError, fmt.Errorf(format, a))
}
