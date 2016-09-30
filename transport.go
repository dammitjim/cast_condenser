package main

import (
	"net/http"
	"strongcast/itunes"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

func searchHandler(
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
