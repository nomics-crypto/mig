package libmig

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

const (
	// DirName is the name of the directory migrations are in
	DirName = "migrations"

	// TableName is the name of the migrations table in the database
	TableName = "migrations"

	// Usage is the command line usage documentation
	Usage = `usage:
  mig <command> [arguments]

commands:
  init                       # Initialize an app with mig
  new <name-of-migration>    # Create a new migration with the given name
  up                         # Run all migrations that haven't been run
  down                       # Run all down migrations for migrations that have been run
  help                       # This usage information`

	// InitMigrationUp is the first migration which creates the migrations table
	InitMigrationUp = `CREATE TABLE migrations (
  version text not null unique,
  ran_at timestamp without time zone default now() not null
);`

	// InitMigrationDown reverses the initial migration
	InitMigrationDown = "DROP TABLE migrations;"
)

var (
	// ErrUsage is returned if mig is used improperly
	ErrUsage = fmt.Errorf(Usage)

	// ErrInvalidMigrationName is returned if a migration name doesn't match ValidMigration
	ErrInvalidMigrationName = fmt.Errorf("invalid migration name: must be letters, numbers, or dashes")

	// ValidMigration is a regexp that only matches valid migration
	// names, which are letters, numbers, and dashes
	ValidMigration = regexp.MustCompile("[a-z0-9-]+")
)

// Run is the entrypoint for the executable. It takes and arguments list
// and returns an error.
func Run(args []string) error {
	if len(args) == 0 {
		return ErrUsage
	}
	switch args[0] {
	case "init":
		return Init()
	case "new":
		if len(args) != 2 {
			return ErrUsage
		}
		return New(args[1])
	case "up":
		return Up()
	case "down":
		return Down()
	}
	return ErrUsage
}

// Init is called via `mig init` and initializes a project for mig.
func Init() error {
	if err := os.Mkdir(DirName, 0755); err != nil {
		return err
	}

	if err := createMigration("create-migrations", InitMigrationUp, InitMigrationDown); err != nil {
		return err
	}

	fmt.Println("mig initialized")

	if err := Up(); err != nil {
		return err
	}

	return nil
}

// New is called via `mig new` and it creates a new migration.
// It expects an argument that is the name of the migration
func New(name string) error {
	return createMigration(name, "/* Your up migration code here */", "/* Your down migration code here */")
}

// Up is called via `mig up` and it runs all migrations that
// have not been run, in order.
func Up() error {
	db, err := db()
	if err != nil {
		return err
	}

	vs, err := versions(db, false)
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(DirName, "*.up.sql"))
	if err != nil {
		return err
	}

	run := []string{}
	for _, f := range files {
		name := filepath.Base(f[:len(f)-7])
		exists := false
		for _, v := range vs {
			if name == v {
				exists = true
				break
			}
		}
		if !exists {
			run = append(run, name)
		}
	}

	for _, r := range run {
		fmt.Println("running", r)
		contents, err := ioutil.ReadFile(filepath.Join(DirName, r+".up.sql"))
		if err != nil {
			return err
		}
		_, err = db.Exec(string(contents))
		if err != nil {
			return err
		}
	}

	if len(run) == 0 {
		fmt.Println("nothing to run")
	}

	return nil
}

// Down runs all down migrations for migrations that have been run
func Down() error {
	db, err := db()
	if err != nil {
		return err
	}

	vs, err := versions(db, true)
	if err != nil {
		return err
	}

	for _, v := range vs {
		fmt.Println("reverting", v)
		contents, err := ioutil.ReadFile(filepath.Join(DirName, v+".down.sql"))
		if err != nil {
			return err
		}
		_, err = db.Exec(string(contents))
		if err != nil {
			return err
		}
	}

	if len(vs) == 0 {
		fmt.Println("nothing to do")
	}

	return nil
}

func createMigration(name, up, down string) error {
	if !ValidMigration.MatchString(name) {
		return ErrInvalidMigrationName
	}
	prefix := date() + "-" + name
	u, err := os.Create(filepath.Join(DirName, prefix+".up.sql"))
	if err != nil {
		return err
	}

	d, err := os.Create(filepath.Join(DirName, prefix+".down.sql"))
	if err != nil {
		return err
	}

	_, err = u.WriteString("BEGIN;\n\n" + up + "\n\nINSERT INTO " + TableName + " (version) VALUES ('" + prefix + "');\nCOMMIT;")
	if err != nil {
		return err
	}

	_, err = d.WriteString("BEGIN;\nDELETE FROM " + TableName + " WHERE version='" + prefix + "';\n\n" + down + "\n\nCOMMIT;")
	if err != nil {
		return err
	}

	return nil
}

func date() string {
	now := time.Now()
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func db() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		envBytes, err := ioutil.ReadFile(".env")
		if err == nil {
			envLines := strings.Split(string(envBytes), "\n")
			for _, line := range envLines {
				tuple := strings.Split(line, "=")
				if len(tuple) > 1 && tuple[0] == "DATABASE_URL" {
					dbURL = strings.Join(tuple[1:], "=")
					break
				}
			}
		}
	}
	return sql.Open("postgres", dbURL)
}

func versions(db *sql.DB, reverse bool) ([]string, error) {
	row := db.QueryRow("SELECT true FROM information_schema.tables WHERE table_name=$1", TableName)
	var check bool
	if err := row.Scan(&check); err == sql.ErrNoRows {
		// No versions in the db, run all files because first will make the table
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}

	vs := []string{}

	order := "ASC"
	if reverse {
		order = "DESC"
	}
	rows, err := db.Query("SELECT version FROM " + TableName + " ORDER BY version " + order)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var v string
		err := rows.Scan(&v)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}
