// app.go

package main

import (
    "context"
    "time"
    "database/sql"
    "log"

    "net/http"
    "strconv"
    "encoding/json"
    "github.com/rs/cors"
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type App struct {
    Router *mux.Router
    DB     *sql.DB
}

func (a *App) Initialize(databaseUrl string) {

    var err error
    a.DB, err = sql.Open("postgres", databaseUrl)
    if err != nil {
        log.Fatal(err)
    }

    query := `CREATE TABLE IF NOT EXISTS sessions(id SERIAL PRIMARY KEY, name text, time text, created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP)`
    ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelfunc()
    res, err := a.DB.ExecContext(ctx, query)
    if err != nil {
        log.Printf("Error %s when creating table", err)
        return
    }
    log.Printf("%s table created", res)

    a.Router = mux.NewRouter()
    a.initializeRoutes()
}

func (a *App) Run(addr string) {
    handler := cors.Default().Handler(a.Router)
    log.Fatal(http.ListenAndServe(":8010", handler))
}

func (a *App) getSession(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid session ID")
        return
    }

    s := session{ID: id}
    if err := s.getSession(a.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "Session not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    respondWithJSON(w, http.StatusOK, s)
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
    count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))

    if count > 50 || count < 1 {
        count = 50
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
    var s session
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&s); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := s.createSession(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, s)
}

func (a *App) updateSession(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid session ID")
        return
    }

    var s session
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&s); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
        return
    }
    defer r.Body.Close()
    s.ID = id

    if err := s.updateSession(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, s)
}

func (a *App) deleteSession(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid Session ID")
        return
    }

    s := session{ID: id}
    if err := s.deleteSession(a.DB); err != nil {
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
