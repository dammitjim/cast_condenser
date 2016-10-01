package api

import (
	"condenser/api/external/itunes"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

// SearchHandler searches both the itunes and django APIs for results.
func SearchHandler(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	qv := r.URL.Query()
	term := qv.Get("term")
	limit := qv.Get("limit")
	// TODO validate things
	if limit == "" {
		limit = "1"
	}

	err := itunes.Search(term, limit)
	if err != nil {
		logrus.Error(err)
	}
}
