package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed migrations/*.sql
	migrationFS embed.FS
)

type DB struct {
	db *sql.DB
	Path string
}

func NewDB(dbPath string) *DB {
	return &DB{
		Path: dbPath,
	}
}

func (db *DB) Open() (err error) {
	conn, err := sql.Open("sqlite3", db.Path)
	if err != nil {
		return err
	}
	db.db = conn

	if _, err = db.db.Exec("pragma journal_mode = wal;"); err != nil {
		return fmt.Errorf("journal_mode failure %v", err)
	}

	if _, err = db.db.Exec("pragma foreign_key = on;"); err != nil {
		return fmt.Errorf("foreign_key failed: %v", err);
	}

	if err = db.runMigrations(); err != nil {
		return fmt.Errorf("migrations failed: %v", err)
	}

	return nil;
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) runMigrations() error {
	db.db.Exec(`create table if not exists migrations (
		filename text primary key
	)`)

	if names, err := fs.Glob(migrationFS, "migrations/*.sql"); err == nil {
		sort.Strings(names)
		for i := 0; i < len(names); i++ {
			if err = db.migrateFile(names[i]); err != nil {
				return err;
			}
		}
	}
	return nil;
}

func (db *DB) migrateFile(filename string) (err error) {
	var count uint8
	err = db.db.QueryRow(`select count(*) from migrations where filename = ?`, filename).Scan(&count)
	if count == 1 {
		fmt.Printf("skipped migrating: %v\n", filename)
		return nil;
	}

	bytes, err := fs.ReadFile(migrationFS, filename)
	if err != nil {
		return fmt.Errorf("reading file '"+filename+"': %v", err)
	}
	tx, err := db.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("opening migration transaction: %v", err)
	}
	
	if _, err := tx.Exec(string(bytes)); err != nil {
		return fmt.Errorf("migrating file '"+filename+"': %v", err)
	}

	if _, err = tx.Exec("insert into migrations (filename) values (?)", filename); err != nil {
		return err
	}

	fmt.Printf("sucessfully migrated: %v", filename)

	tx.Commit()
	return nil;
}