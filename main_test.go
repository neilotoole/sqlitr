package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_execute(t *testing.T) {
	const rowCount = 10

	dbFile := filepath.Join("testdata", "example.sqlite")

	records := mustQuery(t, dbFile, false, "SELECT * FROM actor")
	require.Equal(t, rowCount+1, len(records)) // +1 is for header row

	records = mustQuery(t, dbFile, true, "SELECT * FROM actor")
	require.Equal(t, rowCount, len(records))

	ctx := context.Background()
	buf := &bytes.Buffer{}

	err := execute(ctx, buf, []string{dbFile, "INSERT INTO actor (actor_id, first_name, last_name) VALUES(11, 'Kubla', 'Khan')"})
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Rows Affected: 1")
	buf.Reset()

	records = mustQuery(t, dbFile, true, "SELECT * FROM actor")
	assert.Equal(t, rowCount+1, len(records)) // should be an extra row now

	err = execute(ctx, buf, []string{dbFile, "DELETE FROM actor WHERE first_name = ?", "Kubla"})
	require.NoError(t, err)
	require.Contains(t, buf.String(), "Rows Affected: 1")
	buf.Reset()

	records = mustQuery(t, dbFile, true, "SELECT * FROM actor")
	require.Equal(t, rowCount, len(records))
}

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
