package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data any) error {
	return JSONWithHeaders(w, status, data, nil)
}

func JSONWithHeaders(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

type JSONError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func JSONErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &JSONResponseWriter{ResponseWriter: w}
		next.ServeHTTP(ww, r)
	})
}

type JSONResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *JSONResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.Header().Set("Content-Type", "application/json")
	w.ResponseWriter.WriteHeader(status)
}

func (w *JSONResponseWriter) Write(b []byte) (int, error) {
	if w.status >= 400 {
		// Convert error response to JSON
		jsonError := JSONError{
			Status:  w.status,
			Message: string(b),
		}
		w.ResponseWriter.Header().Set("Content-Type", "application/json")
		return 0, json.NewEncoder(w.ResponseWriter).Encode(jsonError)
	}
	return w.ResponseWriter.Write(b)
}
