package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	validator "gopkg.in/validator.v2"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// App encapsulates Env, Router and middleware
type App struct {
	Config      *Env
	Router      *mux.Router
	Middlewares *Middleware
}

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

type shortlinkResp struct {
	Shortlink string `json:"shortlink"`
}

func (a *App) initializeRoutes() {
	m := alice.New(a.Middlewares.LoggingHandler, a.Middlewares.RecoverHandler)
	a.Router.Handle("/api/shorten", m.ThenFunc(a.createShortlink)).Methods("POST")
	a.Router.Handle("/api/info", m.ThenFunc(a.getShortlinkInfo)).Methods("GET")
	a.Router.Handle("/{shortlink:[a-zA-Z0-9]{1,11}}", m.ThenFunc(a.redirect)).Methods("GET")
}

// Initialize is initialization of app
func (a *App) Initialize(e *Env) {
	// set log formatter
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Config = e
	a.Router = mux.NewRouter()
	a.Middlewares = &Middleware{}
	a.initializeRoutes()
}

// Run starts listen and server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) createShortlink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, StatusError{http.StatusBadRequest, fmt.Errorf("parse parameters failed %v", r.Body)})
		return
	}
	if err := validator.Validate(req); err != nil {
		respondWithError(w, StatusError{http.StatusBadRequest, fmt.Errorf("validate parameters failed %v", req)})
		return
	}
	defer r.Body.Close()

	s, err := a.Config.S.Shorten(req.URL, req.ExpirationInMinutes)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJSON(w, http.StatusCreated, shortlinkResp{Shortlink: s})
	}
}

func (a *App) getShortlinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")
	d, err := a.Config.S.ShortlinkInfo(s)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJSON(w, http.StatusOK, d)
	}
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u, err := a.Config.S.Unshorten(vars["shortlink"])
	if err != nil {
		respondWithError(w, err)
	} else {
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}

func respondWithError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case Error:
		log.Printf("HTTP %d  - %s", e.Status(), e)
		respondWithJSON(w, e.Status(), e.Error())
	default:
		respondWithJSON(w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
