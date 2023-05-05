package repository

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/soldiii/diplom/internal/model"
)

type AdPostgres struct {
	db *sqlx.DB
}

func NewAdPostgres(db *sqlx.DB) *AdPostgres {
	return &AdPostgres{db: db}
}

func (r *AdPostgres) CreateAd(ad *model.Advertisement) (int, error) {
	var id int
	sup_id, err := strconv.Atoi(ad.SupervisorID)
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (supervisor_id, title, text) VALUES ($1, $2, $3) RETURNING id", adsTable)
	row := r.db.QueryRow(query, sup_id, ad.Title, ad.Text)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AdPostgres) GetAdsBySupervisorID(supervisorID string) ([]*model.Advertisement, error) {
	var ads []*model.Advertisement
	query := fmt.Sprintf("SELECT * FROM %s WHERE supervisor_id = $1", adsTable)
	rows, err := r.db.Query(query, supervisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		ad := model.Advertisement{}
		err := rows.Scan(&ad.ID, &ad.SupervisorID, &ad.Title, &ad.Text)
		if err != nil {
			return nil, err
		}
		ads = append(ads, &ad)
	}
	return ads, nil
}

func (r *AdPostgres) UpdateAd(title string, text string, adID string) (int, error) {
	query := fmt.Sprintf("UPDATE %s SET title = $2, text = $3 WHERE id = $1", adsTable)
	r.db.QueryRow(query, adID, title, text)
	id, err := strconv.Atoi(adID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AdPostgres) DeleteAd(adID string) (int, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", adsTable)
	r.db.QueryRow(query, adID)
	id, err := strconv.Atoi(adID)
	if err != nil {
		return 0, err
	}
	return id, nil
}
