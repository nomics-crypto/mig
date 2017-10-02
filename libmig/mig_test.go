package libmig

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestHelp(t *testing.T) {
	err := Run([]string{"help"})
	if err != ErrUsage {
		t.Fatal(err)
	}
}

func TestNoArgs(t *testing.T) {
	err := Run([]string{})
	if err != ErrUsage {
		t.Fatal(err)
	}
}

func TestInitialize(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS migrations")
	if err != nil {
		t.Fatal(err)
	}

	dir, err := ioutil.TempDir("", "mig-test")
	if err != nil {
		t.Fatal(err)
	}
	os.Chdir(dir)

	migrationsDirectory := filepath.Join(dir, "migrations")
	f, err := os.Open(migrationsDirectory)
	if err == nil {
		t.Fatal("expected error opening migrations path")
	}
	if err, ok := err.(*os.PathError); !ok {
		t.Fatal("error is not a path error", err)
	}
	if f != nil {
		t.Fatal("file should be nil")
	}

	err = Run([]string{"init"})
	if err != nil {
		t.Fatal(err)
	}

	f, err = os.Open(migrationsDirectory)
	if err != nil {
		t.Fatal("error opening migrations path after init")
	}

	info, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if !info.IsDir() {
		t.Fatal("migrations should be a folder")
	}

	var versionCount int
	err = db.QueryRow("SELECT count(*) FROM migrations").Scan(&versionCount)

	if versionCount != 1 {
		t.Fatal("expected initial version in database")
	}

	err = Run([]string{"new", "do-something"})
	if err != nil {
		t.Fatal(err)
	}

	err = Run([]string{"up"})
	if err != nil {
		t.Fatal(err)
	}

	err = db.QueryRow("SELECT count(*) FROM migrations").Scan(&versionCount)

	if versionCount != 2 {
		t.Fatal("expected new version in database")
	}

	err = Run([]string{"down"})
	if err != nil {
		t.Fatal(err)
	}

	var check bool
	row := db.QueryRow("SELECT true from information_schema.tables WHERE table_name='migrations'")
	err = row.Scan(&check)
	if err != sql.ErrNoRows {
		t.Fatal("expected no migrations table", err)
	}

}
