package rest

import (
	"context"
	"encoding/json"
	"log"
	"mart"
	"net/http"
	"strconv"
)

// jsType is a http middleware which sets the Content-Type response header.
func jsType(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}

// writeError responses with JSON-encoded data of msg to client by given
// status code.
func writeError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func marts(w http.ResponseWriter, _ *http.Request) {
	var ms []mart.MartInfo
	for _, m := range mart.All() {
		ms = append(ms, m.Info())
	}

	if len(ms) == 0 {
		writeError(w, http.StatusNotFound, "no mart available")
		return
	}

	json.NewEncoder(w).Encode(ms)
}

func search(w http.ResponseWriter, r *http.Request) {

	// parse key, num & order
	key := r.FormValue("key")
	if key == "" {
		writeError(w, http.StatusBadRequest, "key must not be empty")
		return
	}
	num, err := strconv.Atoi(r.FormValue("num"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "num must be numbers")
		return
	}
	odr, _ := strconv.Atoi(r.FormValue("order"))

	// get marts
	var ms []*mart.Mart
	if id := r.FormValue("mart"); id != "" {
		m, err := mart.Open(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "mart "+id+" not available")
			return
		}
		ms = append(ms, m)
	} else {
		ms = mart.All()
		if len(ms) == 0 {
			writeError(w, http.StatusNotFound, "no mart available")
			return
		}
	}

	// create query and make search request
	d := make(chan bool)
	q := mart.Query{
		Key:   key,
		Order: odr,
		Done:  func() { d <- true },
	}

	put := make(chan []mart.Product)
	che := make(chan error)
	ctx, quit := context.WithCancel(context.Background())
	defer quit()

	for i := range ms {
		ms[i].Search(ctx, q, put, che)
	}

	var done int
	var p []mart.Product
	for {
		select {
		case ps := <-put:
			p = append(p, ps...)
			if num > 0 && len(p) > num {
				quit()
				p = p[:num]
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(p)
				return
			}
		case err = <-che:
			log.Println(err)
		case <-d:
			done++
			if done == len(ms) { // all jobs are done
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(p)
				return
			}
		}
	}
}

// Serve creates a RESTful server which listens to given port.
func Serve(port string) error {
	http.HandleFunc("/marts", marts)
	http.HandleFunc("/search", search)
	return http.ListenAndServe(port, nil)
}
