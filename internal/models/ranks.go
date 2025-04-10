package models

import (
	"database/sql"
)

type Rank struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type RankInput struct {
	Name string `json:"name"`
}

// ListRanks retrieves all ranks
func ListRanks(db *sql.DB) ([]Rank, error) {
	rows, err := db.Query(`SELECT id, name FROM rank ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ranks []Rank
	for rows.Next() {
		var r Rank
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}
		ranks = append(ranks, r)
	}
	return ranks, rows.Err()
}

// GetRankByID retrieves a single rank by ID
func GetRankByID(db *sql.DB, id int64) (*Rank, error) {
	var r Rank
	row := db.QueryRow(`SELECT id, name FROM rank WHERE id = $1`, id)
	if err := row.Scan(&r.ID, &r.Name); err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateRank inserts a new rank
func CreateRank(db *sql.DB, r RankInput) (*int64, error) {
	var id int64
	if err := db.QueryRow(
		`INSERT INTO rank (name) VALUES ($1) RETURNING id`,
		r.Name,
	).Scan(&id); err != nil {
		return nil, err
	}
	return &id, nil
}

// UpdateRank updates an existing rank
func UpdateRank(db *sql.DB, rankID int64, r RankInput) (*Rank, error) {
	res, err := db.Exec(`UPDATE rank SET name = $1 WHERE id = $2`, r.Name, rankID)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	updated := &Rank{
		ID:   rankID,
		Name: r.Name,
	}

	return updated, err
}

// DeleteRank deletes a rank by ID
func DeleteRank(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM rank WHERE id = $1`, id)
	return err
}
