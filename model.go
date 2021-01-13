// model.go

package main

import (
    "database/sql"
)

type session struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Time float64 `json:"time"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

func (p *session) getSession(db *sql.DB) error {
  return db.QueryRow("SELECT name, time FROM sessions WHERE id=$1",
      p.ID).Scan(&p.Name, &p.Time)
}

func (p *session) updateSession(db *sql.DB) error {
  _, err :=
      db.Exec("UPDATE sessions SET name=$1, time=$2 WHERE id=$3",
          p.Name, p.Time, p.ID)

  return err
}

func (p *session) deleteSession(db *sql.DB) error {
  _, err := db.Exec("DELETE FROM sessions WHERE id=$1", p.ID)

  return err
}

func (p *session) createSession(db *sql.DB) error {
  err := db.QueryRow(
      "INSERT INTO sessions(name, time) VALUES($1, $2) RETURNING id",
      p.Name, p.Time).Scan(&p.ID)

  if err != nil {
      return err
  }

  return nil
}

func getSessions(db *sql.DB, start, count int) ([]session, error) {
  rows, err := db.Query(
      "SELECT id, name, time, created_at, updated_at FROM sessions LIMIT $1 OFFSET $2",
      count, start)

  if err != nil {
      return nil, err
  }

  defer rows.Close()

  sessions := []session{}

  for rows.Next() {
      var p session
      if err := rows.Scan(&p.ID, &p.Name, &p.Time, &p.CreatedAt, &p.UpdatedAt); err != nil {
          return nil, err
      }
      sessions = append(sessions, p)
  }

  return sessions, nil
}
