package main

import (
	log "log"
	http "net/http"
	os "os"

	_ "github.com/lib/pq"

	db "api/internal/db"
	server "api/internal/server"
)

// main function
func main() {
    // open the connection (and validate with Ping)
    store, err := db.Open(os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatalf("db open: %v", err)
    }
    defer store.DB.Close()


    // run any schema migrations
    if err:= store.Migrate(); err != nil {
        log.Fatalf("db migrate: %v", err)
    }

    // build and run the HTTP router, passing the store
    router := server.NewRouter(store)
    log.Fatal(http.ListenAndServe(":8000", router))
}
