package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// expectRowCount is the number of rows in the "actor" table
	// in testdata/example.sqlite
	expectRowCount = 10

	flagNoHeaderTrue = "--no-header=true"
)

func Test_doQuery_doExec(t *testing.T) {
	dbFile, ctx, buf := testDB(t), context.Background(), &bytes.Buffer{}
	cfg := config{out: buf, dbFile: dbFile, noHeader: true}

	records := doQueryRecords(t, dbFile, true, "SELECT * FROM actor")
	require.Equal(t, expectRowCount, len(records))

	records = doQueryRecords(t, dbFile, false, "SELECT * FROM actor")
	require.Equal(t, expectRowCount+1, len(records)) // +1 is for header row

	err := doExec(ctx, cfg, "INSERT INTO actor (actor_id, first_name, last_name) VALUES(11, 'Kubla', 'Khan')")
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Rows Affected: 1")
	buf.Reset()

	records = doQueryRecords(t, dbFile, true, "SELECT * FROM actor")
	require.Equal(t, expectRowCount+1, len(records)) // should be an extra row now

	err = doExec(ctx, cfg, "DELETE FROM actor WHERE first_name = ?", "Kubla")
	require.NoError(t, err)
	require.Contains(t, buf.String(), "Rows Affected: 1")
	buf.Reset()

	records = doQueryRecords(t, dbFile, true, "SELECT * FROM actor")
	require.Equal(t, expectRowCount, len(records))
}

func Test_cli_exec(t *testing.T) {
	dbFile, ctx, buf := testDB(t), context.Background(), &bytes.Buffer{}
	osArgs := []string{t.Name(), "--exec", dbFile, "DELETE FROM actor WHERE actor_id <= ?", "5"}

	err := cli(ctx, buf, osArgs)
	require.NoError(t, err)
	require.Contains(t, buf.String(), "Rows Affected: 5")
}

func Test_cli_query(t *testing.T) {
	dbFile, ctx, buf := testDB(t), context.Background(), &bytes.Buffer{}
	osArgs := []string{t.Name(), flagNoHeaderTrue, dbFile, "SELECT * FROM actor"}

	err := cli(ctx, buf, osArgs)
	require.NoError(t, err)

	records := readTSV(t, buf)
	require.Equal(t, expectRowCount, len(records))
}

func Test_cli_query_download(t *testing.T) {
	const dbURL = "https://github.com/neilotoole/sqlitr/raw/dev/testdata/example.sqlite"

	ctx, buf := context.Background(), &bytes.Buffer{}
	osArgs := []string{t.Name(), flagNoHeaderTrue, dbURL, "SELECT * FROM actor"}

	err := cli(ctx, buf, osArgs)
	require.NoError(t, err)

	records := readTSV(t, buf)
	require.Equal(t, expectRowCount, len(records))
}

func Test_cli_create(t *testing.T) {
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

func Test_cli_help(t *testing.T) {
	ctx, buf := context.Background(), &bytes.Buffer{}
	osArgs := []string{os.Args[0], "--help"}

	err := cli(ctx, buf, osArgs)
	require.NoError(t, err)

	got := buf.String()
	require.Equal(t, msgHelp, got)
}

func Test_cli_version(t *testing.T) {
	ctx, buf := context.Background(), &bytes.Buffer{}
	osArgs := []string{os.Args[0], "--version"}

	err := cli(ctx, buf, osArgs)
	require.NoError(t, err)

	got := buf.String()
	require.True(t, strings.HasPrefix(got, "sqlitr "))
}

func Test_download(t *testing.T) {
	const wantURL = "https://github.com/neilotoole/sqlitr/raw/master/testdata/example.sqlite"

	ctx := context.Background()
	destDir, err := ioutil.TempDir("", t.Name())
	require.NoError(t, err)

	gotFile, written, err := download(ctx, wantURL, destDir, "")
	require.NoError(t, err)
	t.Logf("downloaded %s  -->  %s", wantURL, gotFile)
	require.True(t, written > 0)
	fi, err := os.Stat(gotFile)
	require.NoError(t, err)
	require.Equal(t, written, fi.Size())
}

// doQueryRecords is a testing convenience method that
// executes query and returns the resulting records (as parsed
// from doQuery's TSV output), failing the test on any error.
func doQueryRecords(t *testing.T, dbFile string, noHeader bool, query string, queryArgs ...interface{}) [][]string {
	ctx, buf := context.Background(), &bytes.Buffer{}
	cfg := config{out: buf, dbFile: dbFile, noHeader: noHeader}

	err := doQuery(ctx, cfg, query, queryArgs...)
	require.NoError(t, err)
	return readTSV(t, buf)
}

func readTSV(t *testing.T, r io.Reader) [][]string {
	csvReader := csv.NewReader(r)
	csvReader.Comma = '\t'
	records, err := csvReader.ReadAll()
	require.NoError(t, err)
	return records
}

// testDB returns the path of a copy of the fixture test database.
// The caller is free to mutate the DB.
func testDB(t *testing.T) string {
	dbFile, cleanup := dbCopy(t, filepath.Join("testdata", "example.sqlite"))
	t.Cleanup(cleanup)
	return dbFile
}

// dbCopy makes a copy of the db file so that git doesn't get annoyed by
// tests touching the version-controlled db file. The path to the
// file copy and a cleanup func to remove the copy are returned.
func dbCopy(t *testing.T, dbFile string) (dbFile2 string, cleanup func()) {
	dbData, err := ioutil.ReadFile(dbFile)
	require.NoError(t, err)
	f, err := ioutil.TempFile("", filepath.Base(dbFile)+"_*")
	require.NoError(t, err)

	dbFile2 = f.Name()
	_, err = f.Write(dbData)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return dbFile2, func() {
		assert.NoError(t, os.Remove(dbFile2))
	}
}
