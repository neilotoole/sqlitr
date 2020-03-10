package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_doQuery_doExec(t *testing.T) {
	const rowCount = 10
	dbFile, cleanup := dbCopy(t, filepath.Join("testdata", "example.sqlite"))
	t.Cleanup(cleanup)

	records := mustQueryRecords(t, dbFile, true, "SELECT * FROM actor")
	require.Equal(t, rowCount, len(records))
	records = mustQueryRecords(t, dbFile, false, "SELECT * FROM actor")
	require.Equal(t, rowCount+1, len(records)) // +1 is for header row

	ctx, buf := context.Background(), &bytes.Buffer{}
	cfg := config{out: buf, dbFile: dbFile}
	err := doExec(ctx, cfg, "INSERT INTO actor (actor_id, first_name, last_name) VALUES(11, 'Kubla', 'Khan')")
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Rows Affected: 1")
	buf.Reset()

	records = mustQueryRecords(t, dbFile, true, "SELECT * FROM actor")
	assert.Equal(t, rowCount+1, len(records)) // should be an extra row now
	// ^ we assert here because we want the DELETE to occur regardless.

	err = doExec(ctx, cfg, "DELETE FROM actor WHERE first_name = ?", "Kubla")
	require.NoError(t, err)
	require.Contains(t, buf.String(), "Rows Affected: 1")
	buf.Reset()

	records = mustQueryRecords(t, dbFile, true, "SELECT * FROM actor")
	require.Equal(t, rowCount, len(records))
}

func Test_doCreate(t *testing.T) {
	ctx, buf := context.Background(), &bytes.Buffer{}

	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, os.RemoveAll(tmpDir)) })

	dbFile := filepath.Join(tmpDir, "sqlitr.db")

	// verify that the file doesn't exist
	_, err = os.Stat(dbFile)
	require.Error(t, err)

	osArgs := []string{os.Args[0], "--create=" + dbFile}
	err = cli(ctx, buf, osArgs)
	require.NoError(t, err)

	got := buf.String()
	want := "Created SQLite DB: " + dbFile + "\n"
	require.Equal(t, want, got)

	// verify that we can open/ping the db
	db, err := sql.Open("sqlite3", dbFile)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })
	require.NoError(t, db.PingContext(ctx))
}

func Test_help(t *testing.T) {
	ctx, buf := context.Background(), &bytes.Buffer{}
	osArgs := []string{os.Args[0], "--help"}

	err := cli(ctx, buf, osArgs)
	require.NoError(t, err)

	got := buf.String()
	require.Equal(t, msgHelp, got)
}

func Test_version(t *testing.T) {
	ctx, buf := context.Background(), &bytes.Buffer{}
	osArgs := []string{os.Args[0], "--version"}

	err := cli(ctx, buf, osArgs)
	require.NoError(t, err)

	got := buf.String()
	require.True(t, strings.HasPrefix(got, "sqlitr "))
}

// mustQueryRecords is a testing convenience method that
// executes query and returns the resulting records (as parsed
// from doQuery's TSV output), failing the test on any error.
func mustQueryRecords(t *testing.T, dbFile string, noHeader bool, query string, queryArgs ...interface{}) [][]string {
	ctx := context.Background()
	buf := &bytes.Buffer{}
	csvReader := csv.NewReader(buf)
	csvReader.Comma = '\t'
	cfg := config{out: buf, dbFile: dbFile, noHeader: noHeader}

	err := doQuery(ctx, cfg, query, queryArgs...)
	require.NoError(t, err)

	records, err := csvReader.ReadAll()
	require.NoError(t, err)

	return records
}

// dbCopy makes a copy of the db file so that git doesn't get annoyed by
// tests touching the version-controlled db file. The path to the
// file copy and a cleanup func to remove the copy are returned.
func dbCopy(t *testing.T, dbFile string) (dbFile2 string, cleanup func()) {
	dbData, err := ioutil.ReadFile(dbFile)
	require.NoError(t, err)
	f, err := ioutil.TempFile("", "*_sqlitr.example.sqlite")
	require.NoError(t, err)

	dbFile2 = f.Name()
	_, err = f.Write(dbData)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return dbFile2, func() {
		assert.NoError(t, os.Remove(dbFile2))
	}
}
