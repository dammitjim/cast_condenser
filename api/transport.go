package api

import (
	"condenser/api/apierrors"
	"condenser/api/external/itunes"
	"condenser/api/process"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func handleError(w http.ResponseWriter, err error) {
	var ae *apierrors.APIError
	switch err.(type) {
	case *apierrors.APIError:
		ae = err.(*apierrors.APIError)
	case apierrors.APIError:
		e := err.(apierrors.APIError)
		ae = &e
	default:
		ae = apierrors.Generic.WithDetails(err.Error())
	}

	http.Error(w, ae.Error(), ae.HTTPStatusCode)
}

// SearchHandler searches both the itunes and django APIs for results.
func SearchHandler(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	qv := r.URL.Query()
	term := qv.Get("term")
	limit := qv.Get("limit")
	if limit == "" {
		limit = "5"
	}

	content, err := itunes.Search(term, limit)
	if err != nil {
		handleError(w, err)
		return
	}

	if len(content.Results) > 0 {
		go process.Run(content.Results...)
	}

	byt, err := json.Marshal(content)
	if err != nil {
		handleError(w, err)
		return
	}

	h := w.Header()
	w.WriteHeader(200)
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "GET")
	h.Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write([]byte(byt))
}
