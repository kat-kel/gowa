package db

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func setupMock(t *testing.T) (*Store, sqlmock.Sqlmock) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    s := &Store{DB: db}
    return s, mock
}

func TestGetActors(t *testing.T) {
    store, mock := setupMock(t)
    defer store.DB.Close()

    rows := sqlmock.NewRows([]string{"id", "name", "voice"}).
        AddRow(1, "Alice", "ALICE").
        AddRow(2, "Bob", "BOB")

    mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, voice FROM actors")).
        WillReturnRows(rows)

    actors, err := store.GetActors(context.Background())
    require.NoError(t, err)
    require.Len(t, actors, 2)
    require.Equal(t, "Alice", actors[0].Name)
    require.Equal(t, "Bob", actors[1].Name)
    require.NoError(t, mock.ExpectationsWereMet())
}

// verify that an empty result set produces a non-nil, zero-length slice
func TestGetActors_Empty(t *testing.T) {
    store, mock := setupMock(t)
    defer store.DB.Close()

    rows := sqlmock.NewRows([]string{"id", "name", "voice"})
    mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, voice FROM actors")).
        WillReturnRows(rows)

    actors, err := store.GetActors(context.Background())
    require.NoError(t, err)
    require.Len(t, actors, 0)
    require.NotNil(t, actors, "slice should be non-nil to encode as [] rather than null")
    require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetActor_NotFound(t *testing.T) {
    store, mock := setupMock(t)
    defer store.DB.Close()

    mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, voice FROM actors WHERE id = $1")).
        WithArgs(42).
        WillReturnError(sql.ErrNoRows)

    _, err := store.GetActor(context.Background(), 42)
    require.ErrorIs(t, err, sql.ErrNoRows)
    require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUpdateDeleteActor(t *testing.T) {
    store, mock := setupMock(t)
    defer store.DB.Close()

    // Create
    a := &Actor{Name: "Carol", Voice: "CAROL"}
    mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO actors (name, voice) VALUES ($1, $2) RETURNING id")).
        WithArgs(a.Name, a.Voice).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(7))
    require.NoError(t, store.CreateActor(context.Background(), a))
    require.Equal(t, 7, a.ID)

    // Update success
    mock.ExpectExec(regexp.QuoteMeta("UPDATE actors SET name=$1, voice=$2 WHERE id=$3")).
        WithArgs(a.Name, a.Voice, a.ID).
        WillReturnResult(sqlmock.NewResult(0, 1))
    require.NoError(t, store.UpdateActor(context.Background(), a))

    // Update missing
    mock.ExpectExec(regexp.QuoteMeta("UPDATE actors SET name=$1, voice=$2 WHERE id=$3")).
        WithArgs(a.Name, a.Voice, a.ID).
        WillReturnResult(sqlmock.NewResult(0, 0))
    err := store.UpdateActor(context.Background(), a)
    require.ErrorIs(t, err, sql.ErrNoRows)

    // Delete success
    mock.ExpectExec(regexp.QuoteMeta("DELETE FROM actors WHERE id=$1")).
        WithArgs(a.ID).
        WillReturnResult(sqlmock.NewResult(0, 1))
    require.NoError(t, store.DeleteActor(context.Background(), a.ID))

    // Delete missing
    mock.ExpectExec(regexp.QuoteMeta("DELETE FROM actors WHERE id=$1")).
        WithArgs(a.ID).
        WillReturnResult(sqlmock.NewResult(0, 0))
    err = store.DeleteActor(context.Background(), a.ID)
    require.ErrorIs(t, err, sql.ErrNoRows)

    require.NoError(t, mock.ExpectationsWereMet())
}
