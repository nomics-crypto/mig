package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestHelp(t *testing.T) {
	err := Run([]string{"help"})
	if err != UsageError {
		t.Fatal(err)
	}
}

func TestNoArgs(t *testing.T) {
	err := Run([]string{})
	if err != UsageError {
		t.Fatal(err)
	}
}

func TestInitialize(t *testing.T) {
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
}
