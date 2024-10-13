package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"new_project/internal/response"
	"runtime/debug"
	"strings"
)

func (app *Server) reportServerError(r *http.Request, err error) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)

	requestAttrs := slog.Group("request", "method", method, "url", url)
	app.logger.Error(message, requestAttrs, "trace", trace)
}

func (app *Server) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]

	err := response.JSONWithHeaders(w, status, map[string]string{"Error": message}, headers)
	if err != nil {
		app.reportServerError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *Server) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.reportServerError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func (app *Server) notFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorMessage(w, r, http.StatusNotFound, message, nil)
}

func (app *Server) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

func (app *Server) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}

//func (app *Server) failedValidation(w http.ResponseWriter, r *http.Request, v validator.Validator) {
//	err := response.JSON(w, http.StatusUnprocessableEntity, v)
//	if err != nil {
//		app.serverError(w, r, err)
//	}
//}
