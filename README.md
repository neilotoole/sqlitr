# sqlitr
`sqlitr` is a trivial query tool for SQLite. It exists solely as a
demonstration for [neilotoole/xcgo](https://github.com/neilotoole/xcgo),
which is a Go cross-compiling docker builder image. `sqlitr` invokes
the SQLite C library via CGo. Building and distributing binaries for
multiple platforms when using CGo is a challenge. This project demonstrates
the use of `neilotoole/xcgo` to make life a bit easier.

## Usage

From `sqlitr --help`:

```
sqlitr is a trivial demonstration query tool for SQLite.

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
```

Running `sqlitr ./testdata/example.sqlite 'SELECT * FROM actor'`:

```tsv
sqlitr ./testdata/example.sqlite 'SELECT * FROM actor'
actor_id	first_name	last_name
1	PENELOPE	GUINESS
2	NICK	WAHLBERG
3	ED	CHASE
4	JENNIFER	DAVIS
5	JOHNNY	LOLLOBRIGIDA
6	BETTE	NICHOLSON
7	GRACE	MOSTEL
8	MATTHEW	JOHANSSON
9	JOE	SWANK
10	CHRISTIAN	GABLE
```

Note that the output is in TSV (tab-separated) format.

## Installation
This section is the raison d'être of `sqlitr`.

The usual Go method will work: `go get -u github.com/neilotoole/sqlitr`


## Acknowledgements
The `testdata/example.sqlite` SQLite database is a tiny
stripped-down version of the [Sakila DB](https://dev.mysql.com/doc/sakila/en/)
with just `10` rows in only one table (`actor`).

`sqlitr` employs [mattn/sqlite3](https://github.com/mattn/sqlite3) to demonstrate
CGo usage.
