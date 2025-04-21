package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// responseWriter is a custom http.ResponseWriter that captures the response data
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	buffer     *bytes.Buffer
}

// WriteHeader captures the status code and passes it to the underlying ResponseWriter
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the response body and passes it to the underlying ResponseWriter
func (rw *responseWriter) Write(b []byte) (int, error) {
	// Write to our buffer first
	rw.buffer.Write(b)
	// Then write to the original response writer
	return rw.ResponseWriter.Write(b)
}

// PrettyPrintJSON formats JSON with indentation
func PrettyPrintJSON(input []byte) string {
	var jsonObj interface{}
	if err := json.Unmarshal(input, &jsonObj); err == nil {
		indented, err := json.MarshalIndent(jsonObj, "", "    ")
		if err == nil {
			return string(indented)
		}
	}

	return string(input)
}

type LoggingDetails struct {
	Method         string
	URI            string
	Proto          string
	Duration       time.Duration
	QueryParams    string
	RequestBody    []byte
	StatusCode     int
	ResponseBody   []byte
	RequestHeader  http.Header
	ResponseHeader http.Header
}

var wg sync.WaitGroup

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var bodyCopy []byte
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				bodyCopy = bodyBytes
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		responseBuffer := &bytes.Buffer{}
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			buffer:         responseBuffer,
		}

		next.ServeHTTP(rw, r)
		duration := time.Since(start)

		details := LoggingDetails{
			Method:         r.Method,
			URI:            r.RequestURI,
			Proto:          r.Proto,
			Duration:       duration,
			QueryParams:    r.URL.Query().Encode(),
			RequestBody:    bodyCopy,
			StatusCode:     rw.statusCode,
			ResponseBody:   responseBuffer.Bytes(),
			RequestHeader:  r.Header,
			ResponseHeader: rw.Header(),
		}

		wg.Add(1)

		go func(details LoggingDetails) {
			defer wg.Done()

			log.Printf("[INFO] %s %s %s %s", details.Method, details.URI, details.Proto, details.Duration)
			log.Println("Query Params:", details.QueryParams)

			if len(details.RequestBody) > 0 {
				contentType := details.RequestHeader.Get("Content-Type")
				if contentType == "application/json" || bytes.HasPrefix(details.RequestBody, []byte("{")) || bytes.HasPrefix(details.RequestBody, []byte("[")) {
					log.Println("Request Body (JSON):")
					log.Println(PrettyPrintJSON(details.RequestBody))
				} else {
					log.Println("Request Body:", string(details.RequestBody))
				}
			}

			if len(details.ResponseBody) > 0 {
				contentType := details.ResponseHeader.Get("Content-Type")

				if contentType == "application/json" || bytes.HasPrefix(details.ResponseBody, []byte("{")) || bytes.HasPrefix(details.ResponseBody, []byte("[")) {
					if len(details.ResponseBody) > 10000 {
						truncated := details.ResponseBody[:10000]
						log.Printf("%d Response Body (JSON):", details.StatusCode)
						log.Println(PrettyPrintJSON(truncated) + "\n... (truncated)")
					} else {
						log.Printf("%d Response Body (JSON):", details.StatusCode)
						log.Println(PrettyPrintJSON(details.ResponseBody))
					}
				} else {
					responseBodyStr := string(details.ResponseBody)
					if len(responseBodyStr) > 1000 {
						responseBodyStr = responseBodyStr[:1000] + "... (truncated)"
					}
					log.Println("Response Body:", responseBodyStr)
				}
			}

			log.Println("-----------------------------------------------------------------------------------------------")
		}(details)
	})
}

func WaitForLogs() {
	wg.Wait()
}
