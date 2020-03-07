// Package main implements the sqlitr demo CLI for neilotoole/xcgo.
// The program executes a query against a SQLite database.
package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)

		<-stop
		cancelFn()
	}()

	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "help", "-h", "-help", "--help":
			fmt.Print(msgHelp)
			return
		}
	}

	err := execute(ctx, os.Stdout, os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func execute(ctx context.Context, out io.Writer, args []string) error {
	// We support only one flag --no-header and it must be the
	// first arg if it is provided at all.
	var noHeader bool
	if len(args) < 2 {
		return errors.New("invalid args")
	}
	if args[0] == "--no-header" {
		noHeader = true
		// We're done with the no-header flag, slice args to
		// get rid of the flag, and carry on.
		args = args[1:]
	}

	// args[0] is sqlite db file path
	// args[1] is the SQL query
	// any additional args are arguments to the SQL query
	if len(args) < 2 {
		return errors.New("invalid args")
	}

	// Create a []interface{} containing any query args
	var queryArgs []interface{}
	// If there's additional args, we append them to queryArgs
	for i := 2; i < len(args); i++ {
		queryArgs = append(queryArgs, args[i])
	}

	// args[0] is the filename
	db, err := sql.Open("sqlite3", args[0])
	if err != nil {
		return err
	}
	defer db.Close()

	query := strings.TrimSpace(args[1])
	// If it's a SELECT, we use db.QueryContext; otherwise db.ExecContext
	if strings.HasPrefix(strings.ToUpper(query), "SELECT") {
		// It's SELECT
		rows, err := db.QueryContext(ctx, query, queryArgs...)
		if err != nil {
			return err
		}

		colNames, err := rows.Columns()
		if err != nil {
			return err
		}

		dests := make([]interface{}, len(colNames))
		for i := range dests {
			// Just make everything into a string, SQLite will handle it fine
			dests[i] = &sql.NullString{}
		}

		w := csv.NewWriter(out)
		w.Comma = '\t'
		var headerWritten bool

		record := make([]string, len(dests))
		for rows.Next() {
			err = rows.Scan(dests...)
			if err != nil {
				return err
			}

			for i := range dests {
				switch v := dests[i].(type) {
				case nil:
					record[i] = ""
				case driver.Valuer:
					val, err := v.Value()
					if err != nil {
						return err
					}
					record[i] = fmt.Sprintf("%v", val)
				default:
					val := reflect.ValueOf(v).Elem()
					record[i] = fmt.Sprintf("%v", val)
				}
			}

			if !noHeader && !headerWritten {
				err := w.Write(colNames)
				if err != nil {
					return err
				}
				headerWritten = true
			}

			err = w.Write(record)
			if err != nil {
				return err
			}
			w.Flush()
		}

		// Done with SELECT query
		return nil
	}

	// Else it's not a SELECT, it's some other SQL statement
	res, err := db.ExecContext(ctx, query, queryArgs...)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "Rows Affected: %d\n", affected)

	lastInserted, err := res.LastInsertId()
	if err == nil && lastInserted > 0 {
		// We don't care about reporting the err from LastInsertId
		fmt.Fprintf(out, "Last Insert ID: %d\n", lastInserted)
	}

	return nil
}

const msgHelp = `sqlitr is a trivial query tool for SQLite.

Usage: sqlitr path/to/db.sqlite query [args]

Examples:
  sqlitr --help
  sqlitr ./testdata/example.sqlite 'SELECT * FROM actor'
  sqlitr --no-header ./testdata/example.sqlite 'SELECT * FROM actor'
  sqlitr ./testdata/example.sqlite "INSERT INTO actor (actor_id, first_name, last_name) VALUES(11, 'Kubla', 'Khan')"
  sqlitr ./testdata/example.sqlite 'DELETE FROM actor WHERE first_name = ?' Kubla

Note that if the query starts with SELECT, output is in TSV (tab-separated)
format. If it's some other SQL statement, the count of rows affected (and
the last insert ID if applicable) are printed.

sqlitr exists solely as a demonstration for neilotoole/xcgo which
is a Go cross-compiling docker builder image. sqlitr was created
by Neil O'Toole <neilotoole@apache.org> and is released under
the MIT License. It is entirely unsupported and will not be developed
further. See https://github.com/neilotoole/sqlitr for more.
`
