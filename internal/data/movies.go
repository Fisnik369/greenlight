package data

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

// `
type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Title     string    `json:"title"`
	Year      int32     `json:"year"`
	Runtime   int32     `json:"runtime"`
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version"`
}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Create(movie *Movie) error {
	query := `INSERT INTO movies (title, year, runtime, genres) VALUES ($1, $2, $3, $4)
              RETURNING id, created_at, version`

	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)} // getting the data from the movie model

	err := m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version) // populating the database
	if err != nil {
		return err
	}
	return nil
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, created_at, title, year, runtime, genres, version
		      FROM movies
			  WHERE id = $1`
	var movie Movie

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}
	return &movie, err
}

func (m MovieModel) Update(movie *Movie) error {
	query := `UPDATE movies
			  SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
			  WHERE id = $5
			  RETURNING version`

	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	return m.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int64) error {
	query := `DELETE FROM movies
			  WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return ErrRecordNotFound
	}

	return nil
}
