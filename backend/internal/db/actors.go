package db

import (
	"context"
	"database/sql"
	"errors"
)

// Actor models a row in the actors table.
type Actor struct {
    ID    int
    Name  string
    Voice string
}

// Store wraps a *sql.DB and provides typed methods.
// (Store is also defined in db.go; keep them consistent.)

// GetActors returns all actors from the database.
func (s *Store) GetActors(ctx context.Context) ([]Actor, error) {
    rows, err := s.DB.QueryContext(ctx, "SELECT id, name, voice FROM actors")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    actors := []Actor{} // start with empty slice so we never return nil
    for rows.Next() {
        var a Actor
        if err := rows.Scan(&a.ID, &a.Name, &a.Voice); err != nil {
            return nil, err
        }
        actors = append(actors, a)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return actors, nil
}

// GetActor looks up a single actor by id.
func (s *Store) GetActor(ctx context.Context, id int) (Actor, error) {
    var a Actor
    err := s.DB.QueryRowContext(ctx, "SELECT id, name, voice FROM actors WHERE id = $1", id).
        Scan(&a.ID, &a.Name, &a.Voice)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return Actor{}, err
        }
        return Actor{}, err
    }
    return a, nil
}

// CreateActor inserts a new actor and returns it with its ID populated.
func (s *Store) CreateActor(ctx context.Context, a *Actor) error {
    return s.DB.QueryRowContext(ctx,
        "INSERT INTO actors (name, voice) VALUES ($1, $2) RETURNING id", a.Name, a.Voice).
        Scan(&a.ID)
}

// UpdateActor updates an existing actor. It returns sql.ErrNoRows if the
// actor doesn't exist.
func (s *Store) UpdateActor(ctx context.Context, a *Actor) error {
    res, err := s.DB.ExecContext(ctx,
        "UPDATE actors SET name=$1, voice=$2 WHERE id=$3", a.Name, a.Voice, a.ID)
    if err != nil {
        return err
    }
    if n, _ := res.RowsAffected(); n == 0 {
        return sql.ErrNoRows
    }
    return nil
}

// DeleteActor removes the actor by id. Returns sql.ErrNoRows when absent.
func (s *Store) DeleteActor(ctx context.Context, id int) error {
    res, err := s.DB.ExecContext(ctx, "DELETE FROM actors WHERE id=$1", id)
    if err != nil {
        return err
    }
    if n, _ := res.RowsAffected(); n == 0 {
        return sql.ErrNoRows
    }
    return nil
}
