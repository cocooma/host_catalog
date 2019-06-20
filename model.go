// model.go

package main

import (
	"database/sql"
)

type host struct {
	ID    int    `json:"id"`
	Name  string  `json:"name"`
	Ip    string  `json:"ip"`
}

func (p *host) getHost(db *sql.DB) error {
	return db.QueryRow("SELECT name, ip FROM hosts WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Ip)
}

func (p *host) updateHost(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE hosts SET name=$1, ip=$2 WHERE id=$3",
			p.Name, p.Ip, p.ID)

	return err
}

func (p *host) deleteHost(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM hosts WHERE id=$1", p.ID)

	return err
}

func (p *host) createHost(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO hosts(name, ip) VALUES($1, $2) RETURNING id",
		p.Name, p.Ip).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func getHosts(db *sql.DB, start, count int) ([]host, error) {
	rows, err := db.Query(
		"SELECT id, name,  ip FROM hosts LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	hosts := []host{}

	for rows.Next() {
		var p host
		if err := rows.Scan(&p.ID, &p.Name, &p.Ip); err != nil {
			return nil, err
		}
		hosts = append(hosts, p)
	}

	return hosts, nil
}