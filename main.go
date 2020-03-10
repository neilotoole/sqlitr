// Package main implements the sqlitr demo CLI for neilotoole/xcgo.
// The program is a trivial Go/CGo front-end for SQLite.
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
	"path/filepath"
	"reflect"
	"strings"

	// Use pflag instead of stdlib flag to support flags
	// after args in the user input
	"github.com/spf13/pflag"

	_ "github.com/mattn/go-sqlite3"
)

var (
	// version info set via ldflags
	version = "dev"
	commit  = "none"
	date    = "unknown"
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

	err := cli(ctx, os.Stdout, os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

type config struct {
	out      io.Writer
	dbFile   string
	noHeader bool
}

// cli is the sqlitr CLI function.
func cli(ctx context.Context, out io.Writer, osArgs []string) error {
	cfg := config{out: out}

	flags := pflag.NewFlagSet(osArgs[0], pflag.ContinueOnError)
	var flagVersion, flagHelp, flagModeExec bool
	var flagCreate string
	flags.StringVar(&flagCreate, "create", "", "create a new SQLite DB at this path")
	flags.BoolVar(&flagVersion, "version", false, "print version info")
	flags.BoolVar(&flagHelp, "help", false, "print help")
	flags.BoolVar(&cfg.noHeader, "no-header", false, "don't print query results header row")
	flags.BoolVar(&flagModeExec, "exec", false, "execute input as statement rather than query")
	err := flags.Parse(osArgs[1:])
	if err != nil {
		return err
	}

	if flagHelp {
		fmt.Fprint(cfg.out, msgHelp)
		return nil
	}

	if flagVersion {
		if version == "dev" {
			// When built using goreleaser or with ldflags set, we will
			// have valid build info (version, date, commit). When built
			// just using go build, those aren't available, so just do this:
			fmt.Fprintln(cfg.out, "sqlitr dev")
			return nil
		}

		fmt.Fprintf(cfg.out, "sqlitr %s  %s  %s\n", version, date, commit)
		return nil
	}

	if flagCreate != "" {
		var err error
		cfg.dbFile, err = filepath.Abs(flagCreate)
		if err != nil {
			return err
		}
		return doCreate(ctx, cfg)
	}

	// cmdArgs[0] is sqlite db file path
	// cmdArgs[1] is the SQL query
	// Any additional args are arguments to the SQL query
	cmdArgs := pflag.Args()
	if len(cmdArgs) < 2 {
		return errors.New("invalid args")
	}

	cfg.dbFile = cmdArgs[0]
	_, err = os.Stat(cfg.dbFile)
	if err != nil {
		return err
	}

	query := strings.TrimSpace(cmdArgs[1])
	var queryArgs []interface{}
	// If there's additional args, we append them to queryArgs
	for i := 2; i < len(cmdArgs); i++ {
		queryArgs = append(queryArgs, cmdArgs[i])
	}

	if flagModeExec {
		return doExec(ctx, cfg, query, queryArgs...)
	}

	return doQuery(ctx, cfg, query, queryArgs...)
}

// doExec executes the query via db.QueryContext.
func doQuery(ctx context.Context, cfg config, query string, queryArgs ...interface{}) error {
	db, err := sql.Open("sqlite3", cfg.dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

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
		// Just treat everything as a string, SQLite will handle it fine
		dests[i] = &sql.NullString{}
	}

	w := csv.NewWriter(cfg.out)
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

		if !cfg.noHeader && !headerWritten {
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

	return nil
}

// doExec executes the query via db.ExecContext.
func doExec(ctx context.Context, cfg config, query string, queryArgs ...interface{}) error {
	db, err := sql.Open("sqlite3", cfg.dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	// The SQL is a statement such as INSERT
	res, err := db.ExecContext(ctx, query, queryArgs...)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Fprintf(cfg.out, "Rows Affected: %d\n", affected)

	lastInserted, err := res.LastInsertId()
	if err == nil && lastInserted > 0 {
		// We don't care about reporting any err from LastInsertId
		fmt.Fprintf(cfg.out, "Last Insert ID: %d\n", lastInserted)
	}

	return nil
}

// doCreate causes cfg.dbFile to be created, and pinged.
func doCreate(ctx context.Context, cfg config) error {
	db, err := sql.Open("sqlite3", cfg.dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintln(cfg.out, "Created SQLite DB:", cfg.dbFile)
	return nil
}

const msgHelp = `sqlitr is a trivial demonstration query tool for SQLite.

Usage: sqlitr [FLAGS] path/to/db.sqlite query [QUERY_ARGS]

Examples:
  sqlitr --help
  sqlitr --version

  # simple select, will print header row
  sqlitr ./testdata/example.sqlite 'SELECT * FROM actor'

  # same as above, but don't print header row
  sqlitr --no-header ./testdata/example.sqlite 'SELECT * FROM actor'

  # execute INSERT stmt
  sqlitr --exec ./testdata/example.sqlite "INSERT INTO actor (actor_id, first_name, last_name) VALUES(11, 'Kubla', 'Khan')"
  
  # same as above, but supplying query args via the command line
  sqlitr --exec ./testdata/example.sqlite 'DELETE FROM actor WHERE first_name = ?' Kubla
  
  # create a new DB file
  sqlitr --create path/to/db.sqlite


Note that if the SQL is a SELECT or other query, output is
in TSV (tab-separated) format. To execute some other SQL statement
such as INSERT, supply the --exec flag. The count of rows affected
(and the last insert ID if applicable) are printed when --exec is
used.

sqlitr exists as a demonstration project for neilotoole/xcgo which
is a Go cross-compiling docker builder image: sqlitr makes use of
the https://github.com/mattn/sqlite3 package which uses CGo to
incorporate SQLite.

sqlitr was created by Neil O'Toole <neilotoole@apache.org> and is
released under the MIT License. See https://github.com/neilotoole/sqlitr
`
