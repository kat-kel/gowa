package handlers

import (
	db "api/internal/db"
	sql "database/sql"
	json "encoding/json"
	log "log"
	stdhttp "net/http"
	"strconv"

	mux "github.com/gorilla/mux"
)

// GetActors returns a handler that lists all actors.
func GetActors(store *db.Store) stdhttp.HandlerFunc {
    return func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
        actors, err := store.GetActors(r.Context())
        if err != nil {
            log.Printf("store.GetActors: %v", err)
            stdhttp.Error(w, "internal error", stdhttp.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(actors)
    }
}

// GetActor returns a handler that fetches a single actor by ID.
func GetActor(store *db.Store) stdhttp.HandlerFunc {
    return func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        actorID, _ := strconv.Atoi(id)
        u, err := store.GetActor(r.Context(), actorID)
        if err != nil {
            if err == sql.ErrNoRows {
                stdhttp.NotFound(w, r)
                return
            }
            log.Printf("query row actor: %v", err)
            stdhttp.Error(w, "internal error", stdhttp.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(u)
    }
}

// CreateActor returns a handler that inserts a new actor.
func CreateActor(store *db.Store) stdhttp.HandlerFunc {
    return func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
        var u db.Actor
        if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
            stdhttp.Error(w, "bad request", stdhttp.StatusBadRequest)
            return
        }

        if err := store.CreateActor(r.Context(), &u); err != nil {
            log.Printf("store.CreateActor: %v", err)
            stdhttp.Error(w, "internal error", stdhttp.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(u)
    }
}

// UpdateActor returns a handler that updates an existing actor.
func UpdateActor(store *db.Store) stdhttp.HandlerFunc {
    return func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
        var u db.Actor
        if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
            stdhttp.Error(w, "bad request", stdhttp.StatusBadRequest)
            return
        }

        vars := mux.Vars(r)
        id := vars["id"]

        actorID, _ := strconv.Atoi(id)
        u.ID = actorID
        if err := store.UpdateActor(r.Context(), &u); err != nil {
            if err == sql.ErrNoRows {
                stdhttp.NotFound(w, r)
                return
            }
            log.Printf("store.UpdateActor: %v", err)
            stdhttp.Error(w, "internal error", stdhttp.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(u)
    }
}

// DeleteActor returns a handler that removes an actor.
func DeleteActor(store *db.Store) stdhttp.HandlerFunc {
    return func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        // make sure it exists
        actorID, _ := strconv.Atoi(id)
        if err := store.DeleteActor(r.Context(), actorID); err != nil {
            if err == sql.ErrNoRows {
                stdhttp.NotFound(w, r)
                return
            }
            log.Printf("store.DeleteActor: %v", err)
            stdhttp.Error(w, "internal error", stdhttp.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(map[string]string{"message": "actor deleted"})
    }
}
