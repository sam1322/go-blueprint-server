package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
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

// Check if the body is a valid JSON by checking the first few characters
func isJSON(b []byte) bool {
	// Trim leading spaces or newlines and check for basic JSON patterns
	b = []byte(strings.TrimSpace(string(b)))
	return len(b) > 0 && (b[0] == '{' || b[0] == '[')
}

//func JSONErrorMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ww := &JSONResponseWriter{ResponseWriter: w}
//		next.ServeHTTP(ww, r)
//	})
//}

func JSONErrorMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := &JSONResponseWriter{ResponseWriter: w, logger: logger, request: r}
			next.ServeHTTP(ww, r)
		})
	}
}

type JSONResponseWriter struct {
	http.ResponseWriter
	status  int
	logger  *slog.Logger
	request *http.Request
}

func (w *JSONResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.Header().Set("Content-Type", "application/json")
	w.ResponseWriter.WriteHeader(status)
}

func (w *JSONResponseWriter) Write(b []byte) (int, error) {
	if w.status >= 400 {
		if isJSON(b) {
			// If the body is already JSON, don't wrap it again, just send it as is.
			w.ResponseWriter.Header().Set("Content-Type", "application/json")
			return w.ResponseWriter.Write(b)
		}
		// Remove the trailing newline if it exists
		trimmedMessage := strings.TrimRight(string(b), "\n")

		// Convert error response to JSON
		jsonError := JSONError{
			Status:  w.status,
			Message: trimmedMessage,
		}
		w.logger.Error("HTTP error",
			slog.Int("status", w.status),
			slog.Any("message", trimmedMessage),
			slog.String("path", w.request.URL.Path),
			slog.String("method", w.request.Method),
		)

		w.ResponseWriter.Header().Set("Content-Type", "application/json")
		return 0, json.NewEncoder(w.ResponseWriter).Encode(jsonError)
	}
	return w.ResponseWriter.Write(b)
}
