package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

//Get returns one movie and error, if any
func (m *DBModel) Get(id int) (*Movie, error) {

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	movieQuery := `
		SELECT ID, TITLE, DESCRIPTION, YEAR, RELEASE_DATE, RUNTIME,	RATING,
		MPAA_RATING, CREATED_AT, UPDATED_AT 
		FROM MOVIES 
		WHERE ID = $1
	`

	row := m.DB.QueryRowContext(context, movieQuery, id)

	var movie Movie

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.RealeseDate,
		&movie.Runtime,
		&movie.Rating,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	m.movieGenres(&movie)

	return &movie, nil
}

//All returns all movies and error, if any
func (m *DBModel) All() ([]*Movie, error) {

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	moviesQuery := `
		SELECT ID, TITLE, DESCRIPTION, YEAR, RELEASE_DATE, RUNTIME,	RATING,
		MPAA_RATING 
		FROM MOVIES 
		ORDER BY TITLE
	`
	rows, err := m.DB.QueryContext(context, moviesQuery)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movies []*Movie

	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.RealeseDate,
			&movie.Runtime,
			&movie.Rating,
			&movie.MPAARating,
		)

		if err != nil {
			return nil, err
		}

		m.movieGenres(&movie)
		movies = append(movies, &movie)
	}

	return movies, nil
}

//movieGenres gets all the genres of a movie, if any
func (m *DBModel) movieGenres(movie *Movie) error {

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	genreQuery := `
			SELECT MOVIES_GENRES.ID, GENRES.GENRE_NAME 
			FROM MOVIES_GENRES 
			LEFT JOIN GENRES ON (GENRES.ID = MOVIES_GENRES.GENRE_ID)
			WHERE MOVIES_GENRES.MOVIE_ID = $1
		`

	genreRows, _ := m.DB.QueryContext(context, genreQuery, movie.ID)

	defer genreRows.Close()

	genres := make(map[int]string)

	for genreRows.Next() {
		var movieGenre MovieGenre
		err := genreRows.Scan(
			&movieGenre.ID,
			&movieGenre.Genre.GenreName,
		)

		if err != nil {
			return err
		}

		genres[movieGenre.ID] = movieGenre.Genre.GenreName
	}

	movie.Genres = genres

	return nil
}
