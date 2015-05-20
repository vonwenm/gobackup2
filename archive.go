package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

/**
 * archive represents a sql.DB
 */
type archive struct {
	conn *sql.DB
}

/**
 * NewArchive returns a archive instance for a sqlite3 database
 * @param path Path to the sqlite3 file
 * @return db a pointer to a db instance, or nil if an error ocurred
 * @return error An error if something went wrong, nil otherwise
 */
func NewArchive(path string) (*archive, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	archive := &archive{conn: conn}
	if err := archive.ensureTables(); err != nil {
		return nil, err
	}
	return archive, nil
}

/**
 * ensureTables creates the required tables in the
 * database in case they don't yet exist
 */
func (a *archive) ensureTables() error {
	queries := [...]string{
		"CREATE TABLE IF NOT EXISTS file (hash text, filename text, is_deleted boolean, PRIMARY KEY(hash, filename))",
		"CREATE TABLE IF NOT EXISTS upload (hash text, amazon_id text, PRIMARY KEY(hash, amazon_id))",
	}

	for _, query := range queries {
		_, err := a.conn.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

/**
 * Retrieve the sql connection from the archive
 * Used in tests. Do not use otherwise.
 */
func (a *archive) Connection() *sql.DB {
	return a.conn
}

/**
 * AddFile adds a file to the archive
 */
func (a *archive) AddFile(file *ArchivedFile) error {
	stmt, err := a.conn.Prepare("INSERT INTO file(hash, filename, is_deleted) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(file.Hash(), file.Filename(), false)
	if err != nil {
		return err
	}

	stmt, err = a.conn.Prepare("INSERT OR IGNORE INTO upload(hash, amazon_id) VALUES(?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(file.Hash(), file.AmazonId())
	if err != nil {
		return err
	}

	return nil
}

func (a *archive) ListFiles() ([]*ArchivedFile, error) {
	stmt, err := a.conn.Prepare("SELECT f.hash, f.filename, f.is_deleted, u.amazon_id FROM file AS f INNER JOIN upload AS u ON f.hash=u.hash WHERE f.is_deleted=0")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var files []*ArchivedFile
	var hash, filename, amazonId string
	var isDeleted bool
	for rows.Next() {
		rows.Scan(&hash, &filename, &isDeleted, &amazonId)
		files = append(files, &ArchivedFile{
			filename:  filename,
			hash:      hash,
			amazonId:  amazonId,
			isDeleted: isDeleted,
		})
	}

	return files, nil
}

/**
 * FindFileByFilename returns an ArchivedFile given a filename
 */
func (a *archive) FindFileByFilename(filename string) (*ArchivedFile, error) {
	stmt, err := a.conn.Prepare("SELECT f.hash, f.filename, f.is_deleted, u.amazon_id FROM file AS f INNER JOIN upload AS u ON f.hash=u.hash WHERE f.filename=? AND f.is_deleted=0")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var hash, fname, amazonId string
	var isDeleted bool
	err = stmt.QueryRow(filename).Scan(&hash, &fname, &isDeleted, &amazonId)
	if err != nil {
		return nil, err
	}

	return &ArchivedFile{
		filename:  fname,
		hash:      hash,
		amazonId:  amazonId,
		isDeleted: isDeleted,
	}, nil
}

/**
 * FindAmazonIdByHash returns the amazon id of a given hash
 * If there is no such file known the function returns an error
 */
func (a *archive) FindAmazonIdByHash(hash string) (*string, error) {
	stmt, err := a.conn.Prepare("SELECT amazon_id FROM upload WHERE hash=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var amazonId string
	err = stmt.QueryRow(hash).Scan(&amazonId)
	if err != nil {
		return nil, err
	}

	return &amazonId, nil
}

/**
 * DeleteFile marks a file as deleted
 */
func (a *archive) DeleteFile(hash, filename string) error {
	stmt, err := a.conn.Prepare("UPDATE file SET is_deleted=1 WHERE hash=? AND filename=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(hash, filename)
	if err != nil {
		return err
	}

	return nil
}
