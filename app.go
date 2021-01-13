// app.go

package main

import (
    "database/sql"
    "fmt"
    "log"

    "net/http"
    "strconv"
    "encoding/json"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type App struct {
    Router *mux.Router
    DB     *sql.DB
}

func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (a *App) Initialize(user, password, dbname string) {
    connectionString :=
        fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

    var err error
    a.DB, err = sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    a.Router = mux.NewRouter()

    a.initializeRoutes()
}

func (a *App) Run(addr string) {
    log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) getSession(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid session ID")
        return
    }

    p := session{ID: id}
    if err := p.getSession(a.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "Session not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    respondWithJSON(w, http.StatusOK, p)
}


func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func (a *App) getSessions(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))

    if count > 10 || count < 1 {
        count = 10
    }
    if start < 0 {
        start = 0
    }

    sessions, err := getSessions(a.DB, start, count)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, sessions)
}

func (a *App) createSession(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    var p session
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&p); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := p.createSession(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateSession(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid session ID")
        return
    }

    var p session
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&p); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
        return
    }
    defer r.Body.Close()
    p.ID = id

    if err := p.updateSession(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteSession(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid Session ID")
        return
    }

    p := session{ID: id}
    if err := p.deleteSession(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) initializeRoutes() {
    a.Router.HandleFunc("/sessions", a.getSessions).Methods("GET")
    a.Router.HandleFunc("/session", a.createSession).Methods("POST")
    a.Router.HandleFunc("/session/{id:[0-9]+}", a.getSession).Methods("GET")
    a.Router.HandleFunc("/session/{id:[0-9]+}", a.updateSession).Methods("PUT")
    a.Router.HandleFunc("/session/{id:[0-9]+}", a.deleteSession).Methods("DELETE")
}
