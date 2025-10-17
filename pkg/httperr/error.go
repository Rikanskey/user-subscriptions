package httperr

import (
	"github.com/go-chi/render"
	"net/http"
	"user-subscriptions/pkg/logger"
)

func BadRequest(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Bad Request", http.StatusBadRequest)
}

func NotFound(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Not Found", http.StatusNotFound)
}

func UnprocessableEntity(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Unprocessable Entity", http.StatusUnprocessableEntity)
}

func InternalServerError(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Internal Server Error", http.StatusInternalServerError)
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, r *http.Request, logMSg string, status int) {
	logger.GetLogEntry(r).WithError(err).WithField("error-slug", slug).Warn(logMSg)

	var details string
	if err != nil {
		details = err.Error()
	}

	resp := ErrorResponse{slug, details, status}

	if err := render.Render(w, r, resp); err != nil {
		panic(err)
	}
}

type ErrorResponse struct {
	Slug       string `json:"slug"`
	Details    string `json:"details"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(e.httpStatus)

	return nil
}
