package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Reservation struct {
	Id       int
	Start    string
	End      string
	Duration int
}

type ReservationStore struct {
	Db *sql.DB
}

func (rs *ReservationStore) Insert(ctx context.Context, r *Reservation) {
	const query = "INSERT INTO reservations (start_time, end_time, duration) VALUES (?, ?, ?)"
	_, err := rs.Db.ExecContext(ctx, query, r.Start, r.End, r.Duration)
	if err != nil {
		log.Fatal(err)
	}
}

func (rs *ReservationStore) GetById(ctx context.Context, id int) Reservation {
	var res Reservation
	const query = "SELECT id, start_time, end_time, duration FROM reservations WHERE id = ?"
	err := rs.Db.QueryRowContext(ctx, query, id).Scan(&res.Id, &res.Start, &res.End, &res.Duration)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func (rs *ReservationStore) DeleteById(ctx context.Context, id int) {
	var query = "DELETE FROM reservations WHERE id = ?"
	_, err := rs.Db.ExecContext(ctx, query, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Reservation with ID %d deleted.\n", id)
}

func (rs *ReservationStore) Count(ctx context.Context) int {
	var count int
	var query = "SELECT COUNT(*) FROM reservations"
	err := rs.Db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}

func (rs *ReservationStore) GetAll(ctx context.Context) []Reservation {
	var res []Reservation
	var query = "SELECT id, start_time, end_time, duration FROM reservations"
	rows, err := rs.Db.QueryContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var r Reservation
		err := rows.Scan(&r.Id, &r.Start, &r.End, &r.Duration)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, r)
	}
	return res
}
