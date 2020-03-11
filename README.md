# sqlitr
`sqlitr` is a trivial query tool for SQLite. It was created as a
demonstration for [neilotoole/xcgo](https://github.com/neilotoole/xcgo),
which is a Go cross-compiling docker builder image. `sqlitr` invokes
the SQLite C library via CGo... building and distributing binaries for
multiple platforms with CGo is a challenge. The `neilotoole/xcgo`
image makes life easier for the common cases.

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

  # same query, but the SQLite db is first downloaded from
  # the URL to a temp file, then the query is executed. 
  sqlitr https://github.com/neilotoole/sqlitr/raw/master/testdata/example.sqlite 'SELECT * FROM actor'

  # execute a SQL stmt (note the --exec flag, as opposed to default query behavior)
  sqlitr --exec ./testdata/example.sqlite "INSERT INTO actor (actor_id, first_name, last_name) VALUES(11, 'Kubla', 'Khan')"
  
  # execute a SQL stmt, but supply query args via the command line
  sqlitr --exec ./testdata/example.sqlite 'DELETE FROM actor WHERE actor_id = ?' 3
  
  # create a new DB file
  sqlitr --create path/to/db.sqlite


Note that if the SQL is a SELECT or similar query, output is
in TSV (tab-separated) format. To execute a non-query SQL statement
such as INSERT, supply the --exec flag: in that case the count of
rows affected (and the last insert ID if applicable) are printed.

sqlitr exists as a demonstration project for neilotoole/xcgo which
is a Go/CGo cross-compiling docker builder image. sqlitr makes use of
the https://github.com/mattn/sqlite3 package which uses CGo to
incorporate SQLite.

sqlitr was created by Neil O'Toole <neilotoole@apache.org> and is
released under the MIT License. See https://github.com/neilotoole/sqlitr
```

Usage example, with a remote DB file:

```shell script
$ sqlitr https://github.com/neilotoole/sqlitr/raw/master/testdata/example.sqlite 'SELECT * FROM actor'
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


## Installation
This section is the raison d'être of `sqlitr`. After any of these
methods, run `sqlitr --version` to verify your install.

### go get
The usual Go method will work: 

```shell script
go get -u github.com/neilotoole/sqlitr
```

### go install

```shell script
$ git clone https://github.com/neilotoole/sqlitr.git
$ cd sqlitr
$ go install
```

### brew
Install on macos via [brew](https://brew.sh/)

```shell script
TODO

```

### scoop
Install on Windows via [scoop](https://brew.sh/)

```shell script
TODO
```

### snap
Install on Linux via [snap](https://snapcraft.io/docs/getting-started)

```shell script
TODO
```

### Docker
Run directly from the published docker image:

```shell script
docker run neilotoole/sqlitr:latest /example.sqlite 'SELECT * FROM actor'
```

### RPM

### GitHub Release
Download the appropriate file from GitHub [releases], and extract the binary from the archive.


## Acknowledgements
The `testdata/example.sqlite` SQLite database is a tiny
stripped-down version of the [Sakila DB](https://dev.mysql.com/doc/sakila/en/)
with just `10` rows in only one table (`actor`).

`sqlitr` employs [mattn/sqlite3](https://github.com/mattn/sqlite3) to demonstrate
CGo usage.
