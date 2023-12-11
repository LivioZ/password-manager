package vault

import (
	"database/sql"
)

type VaultEntry struct {
	Id       int
	Name     string
	Username string
	Password string
	Other    string
}

func CreateDB(dbFilePath string) error {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	sqlStmt := `
	DROP TABLE IF EXISTS vault;
	CREATE TABLE vault(
	    id INTEGER PRIMARY KEY,
	    name TEXT default '',
	    username TEXT default '',
	    password TEXT default '',
	    other TEXT default '',
	    UNIQUE(name, username));
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func SearchVaultEntries(dbFilePath string, searchTerm string) (result []VaultEntry, err error) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	rows, err := db.Query("SELECT * FROM vault WHERE name LIKE ?", "%"+searchTerm+"%")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		entry := VaultEntry{}
		err = rows.Scan(
			&entry.Id,
			&entry.Name,
			&entry.Username,
			&entry.Password,
			&entry.Other,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, entry)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ListVaultEntries(dbFilePath string) (result []VaultEntry, err error) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	rows, err := db.Query("SELECT * FROM vault")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		entry := VaultEntry{}
		err = rows.Scan(
			&entry.Id,
			&entry.Name,
			&entry.Username,
			&entry.Password,
			&entry.Other,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, entry)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func AddEntry(dbFilePath string, name string, username string, password string, other string) error {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	_, err = db.Exec("INSERT INTO vault(name, username, password, other) VALUES(?, ?, ?, ?)", name, username, password, other)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEntry(dbFilePath string, id int) error {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	_, err = db.Exec("DELETE FROM vault WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}
