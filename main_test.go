package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_execute(t *testing.T) {
	const rowCount = 10

	// make a copy of the db file so that git doesn't get annoyed by
	// this test touching the file
	dbData, err := ioutil.ReadFile(filepath.Join("testdata", "example.sqlite"))
	require.NoError(t, err)
	f, err := ioutil.TempFile("", "*_sqlitr.example.sqlite")
	require.NoError(t, err)
	dbFile := f.Name()
	t.Logf("Using copy of example.sqlite: %s", dbFile)
	require.NoError(t, ioutil.WriteFile(dbFile, dbData, os.ModePerm))
	defer func() { assert.NoError(t, os.Remove(dbFile)) }()

	records := mustQuery(t, dbFile, true, "SELECT * FROM actor")
	require.Equal(t, rowCount, len(records))
	records = mustQuery(t, dbFile, false, "SELECT * FROM actor")
	require.Equal(t, rowCount+1, len(records)) // +1 is for header row

	ctx, buf := context.Background(), &bytes.Buffer{}

	err = execute(ctx, buf, []string{dbFile, "INSERT INTO actor (actor_id, first_name, last_name) VALUES(11, 'Kubla', 'Khan')"})
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Rows Affected: 1")
	buf.Reset()

	records = mustQuery(t, dbFile, true, "SELECT * FROM actor")
	assert.Equal(t, rowCount+1, len(records)) // should be an extra row now
	// ^ we assert here because we want the DELETE to occur regardless.

	err = execute(ctx, buf, []string{dbFile, "DELETE FROM actor WHERE first_name = ?", "Kubla"})
	require.NoError(t, err)
	require.Contains(t, buf.String(), "Rows Affected: 1")
	buf.Reset()

	records = mustQuery(t, dbFile, true, "SELECT * FROM actor")
	require.Equal(t, rowCount, len(records))
}

// mustQuery is a testing convenience method that executes query and
// returns the resulting records, failing the test on any error.
func mustQuery(t *testing.T, dbFile string, noHeader bool, query string) [][]string {
	ctx := context.Background()
	buf := &bytes.Buffer{}
	csvReader := csv.NewReader(buf)
	csvReader.Comma = '\t'

	args := []string{dbFile, query}
	if noHeader {
		args = append([]string{"--no-header"}, args...)
	}

	err := execute(ctx, buf, args)
	require.NoError(t, err)

	records, err := csvReader.ReadAll()
	require.NoError(t, err)

	return records
}
