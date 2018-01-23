package repository

import (
	"database/sql"

	"github.com/src-d/code-annotation/server/model"
)

type Experiments struct {
	db *sql.DB
}

func NewExperiments(db *sql.DB) *Experiments {
	return &Experiments{db}
}

func (r *Experiments) Get(id int) (model.Experiment, error) {
	var exp model.Experiment
	err := r.db.QueryRow("SELECT id, name FROM experiments WHERE id=?", id).
		Scan(&exp.ID, &exp.Name)
	return exp, sqlE(err)
}

func (r *Experiments) HasAssignments(expID int, uID int) (bool, error) {
	var count int
	rows, err := r.db.Query("SELECT COUNT(*) as count FROM assignments WHERE experiment_id=? AND user_id=?", expID, uID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return false, err
		}
	}
	return count != 0, nil
}

func (r *Experiments) CreateAssignments(expID int, uID int) error {
	tx, err := r.db.Begin()
	rows, err := tx.Query("SELECT id FROM file_pairs")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var fpID int
		if err := rows.Scan(&fpID); err != nil {
			tx.Rollback()
			return err
		}
		_, err := tx.Exec(
			"INSERT INTO assignments (experiment_id, pair_id, user_id) VALUES (?, ?, ?)",
			expID, fpID, uID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := rows.Err(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Experiments) GetAssignment(id int) (model.Assignment, error) {
	var a model.Assignment
	err := r.db.QueryRow("SELECT id, pair_id, answer, duration FROM assignments WHERE id=?", id).
		Scan(&a.ID, &a.PairID, &a.Answer, &a.Duration)
	return a, sqlE(err)
}

func (r *Experiments) GetAssignments(expID int, uID int) ([]model.Assignment, error) {
	rows, err := r.db.Query(
		"SELECT id, pair_id, answer, duration FROM assignments WHERE experiment_id=? AND user_id=?",
		expID, uID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var as []model.Assignment
	for rows.Next() {
		var a model.Assignment
		if err := rows.Scan(&a.ID, &a.PairID, &a.Answer, &a.Duration); err != nil {
			return nil, err
		}
		as = append(as, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return as, nil
}

func (r *Experiments) UpdateAssignments(m model.Assignment) error {
	_, err := r.db.Exec(
		"UPDATE assignments SET answer=?, duration=? WHERE id=?",
		m.Answer, m.Duration, m.ID,
	)
	return err
}

func (r *Experiments) GetFilePair(id int) (model.FilePair, error) {
	var fp model.FilePair
	err := r.db.QueryRow("SELECT id, diff FROM file_pairs WHERE id=?", id).
		Scan(&fp.ID, &fp.Diff)
	return fp, sqlE(err)
}
