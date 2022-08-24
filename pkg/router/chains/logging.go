package chains

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	// struct for holding response details
	responseData struct {
		status int
		size   int
	}

	// our http.ResponseWriter implementation
	loggingResponseWriter struct {
		http.ResponseWriter // compose original http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(handler http.Handler, logger *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lrw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		handler.ServeHTTP(&lrw, req)
		logger.Infow("request complete",
			"uri", req.RequestURI,
			"method", req.Method,
			"status", responseData.status,
			"duration", time.Since(start),
			// if request's header does not contain 'Content-Length', then it will be -1
			"req_size", req.ContentLength,
			"resp_size", responseData.size,
			"ua", req.UserAgent(),
			// simplely get from remote addr,
			// but not from x-forward-for header.
			"remoteaddr", req.RemoteAddr,
		)
	})
}
