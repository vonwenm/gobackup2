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
 * AddFile adds a file to the archive
 */
func (a *archive) AddFile(file *ArchivedFile) error {
	// "INSERT INTO file(hash, filename, is_deleted) VALUES (?, ?, ?)"
	// "INSERT OR IGNORE INTO upload(hash, amazon_id) VALUES(?, ?)"

	return nil
}

/**
 * FindFileByFilename returns an ArchivedFile given a filename
 */
func (a *archive) FindFileByFilename(filename string) *ArchivedFile {
	// "SELECT f.hash, f.filename FROM file AS f INNER JOIN upload AS u ON a.hash=u.hash WHERE f.filename=? AND f.is_deleted=0"

	return nil
}

/**
 * FindAmazonIdByHash returns the amazon id of a given hash
 * If there is no such file known the function returns nil
 */
func (a *archive) FindAmazonIdByHash(hash string) *string {
	// "SELECT amazon_id FROM upload WHERE hash=?"

	return nil
}

/**
 * DeleteFile marks a file as deleted
 */
func (a *archive) DeleteFile(hash, filename string) error {
	// "UPDATE file SET is_deleted=1 WHERE hash=? AND filename=?"

	return nil
}
