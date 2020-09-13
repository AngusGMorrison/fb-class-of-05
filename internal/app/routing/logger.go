package routing

import "net/http"

type logger interface {
	Printf(format string, v ...interface{})
}

type loggingHandlerFunc = func(w http.ResponseWriter, r *http.Request, log logger)

type loggingHandler struct {
	logger
	loggingHandlerFunc
}

func (lh *loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lh.loggingHandlerFunc(w, r, lh.logger)
}

func loggingHandlerFactory(log logger) func(hf loggingHandlerFunc) *loggingHandler {
	return func(hf loggingHandlerFunc) *loggingHandler {
		return &loggingHandler{log, hf}
	}
}
