# sqlite3shell
[![Build status](https://ci.appveyor.com/api/projects/status/7ucmixmtbwgh2n1g/branch/master?svg=true)](https://ci.appveyor.com/project/cznic/sqlite3shell/branch/master)

Command sqlite3shell is a mechanically produced Go port of shell.c, part of
the SQLite project.

### Installation

To install and/or update perform

    $ go get [-u] github.com/cznic/sqlite3shell

### Changelog

2018-03-06: Initial release.

### Status

This is an early technology preview of sqlite2go - a support-sqlite-only
version of CCGO[0]. DO NOT USE IN PRODUCTION. Please help the development by
reporting any issues found at [1].

### Supported platforms and operating systems

In GOOS_GOARCH form

    linux_386
    linux_amd64

Windows port is planned and expected to be the next port. Contributors of
other OS/platforms are welcome.

[0] [http://github.com/cznic/ccgo](http://github.com/cznic/ccgo)

[1] [http://github.com/cznic/sqlite3shell/issues](http://github.com/cznic/sqlite3shell/issues)

### Getting help

Usage hints can be obtained via the -help command flag or the .help shell
command.

    $ sqlite3shell -help
    Usage: ./sqlite3shell [OPTIONS] FILENAME [SQL]
    FILENAME is the name of an SQLite database. A new database is created
    if the file does not previously exist.
    OPTIONS include:
       -ascii               set output mode to 'ascii'
       -bail                stop after hitting an error
       -batch               force batch I/O
       -column              set output mode to 'column'
       -cmd COMMAND         run "COMMAND" before reading stdin
       -csv                 set output mode to 'csv'
       -echo                print commands before execution
       -init FILENAME       read/process named file
       -[no]header          turn headers on or off
       -help                show this message
       -html                set output mode to HTML
       -interactive         force interactive I/O
       -line                set output mode to 'line'
       -list                set output mode to 'list'
       -lookaside SIZE N    use N entries of SZ bytes for lookaside memory
       -mmap N              default mmap size set to N
       -newline SEP         set output row separator. Default: '\n'
       -nullvalue TEXT      set text string for NULL values. Default ''
       -pagecache SIZE N    use N slots of SZ bytes each for page cache memory
       -quote               set output mode to 'quote'
       -separator SEP       set output column separator. Default: '|'
       -stats               print memory stats before each finalize
       -version             show SQLite version
       -vfs NAME            use NAME as the default VFS
    $ sqlite3shell db
    SQLite version 3.21.0 2017-10-24 18:55:49
    Enter ".help" for usage hints.
    sqlite> .help
    .auth ON|OFF           Show authorizer callbacks
    .backup ?DB? FILE      Backup DB (default "main") to FILE
    .bail on|off           Stop after hitting an error.  Default OFF
    .binary on|off         Turn binary output on or off.  Default OFF
    .cd DIRECTORY          Change the working directory to DIRECTORY
    .changes on|off        Show number of rows changed by SQL
    .check GLOB            Fail if output since .testcase does not match
    .clone NEWDB           Clone data into NEWDB from the existing database
    .databases             List names and files of attached databases
    .dbinfo ?DB?           Show status information about the database
    .dump ?TABLE? ...      Dump the database in an SQL text format
                             If TABLE specified, only dump tables matching
                             LIKE pattern TABLE.
    .echo on|off           Turn command echo on or off
    .eqp on|off|full       Enable or disable automatic EXPLAIN QUERY PLAN
    .exit                  Exit this program
    .fullschema ?--indent? Show schema and the content of sqlite_stat tables
    .headers on|off        Turn display of headers on or off
    .help                  Show this message
    .import FILE TABLE     Import data from FILE into TABLE
    .imposter INDEX TABLE  Create imposter table TABLE on index INDEX
    .indexes ?TABLE?       Show names of all indexes
                             If TABLE specified, only show indexes for tables
                             matching LIKE pattern TABLE.
    .limit ?LIMIT? ?VAL?   Display or change the value of an SQLITE_LIMIT
    .lint OPTIONS          Report potential schema issues. Options:
                             fkey-indexes     Find missing foreign key indexes
    .load FILE ?ENTRY?     Load an extension library
    .log FILE|off          Turn logging on or off.  FILE can be stderr/stdout
    .mode MODE ?TABLE?     Set output mode where MODE is one of:
                             ascii    Columns/rows delimited by 0x1F and 0x1E
                             csv      Comma-separated values
                             column   Left-aligned columns.  (See .width)
                             html     HTML <table> code
                             insert   SQL insert statements for TABLE
                             line     One value per line
                             list     Values delimited by "|"
                             quote    Escape answers as for SQL
                             tabs     Tab-separated values
                             tcl      TCL list elements
    .nullvalue STRING      Use STRING in place of NULL values
    .once FILENAME         Output for the next SQL command only to FILENAME
    .open ?OPTIONS? ?FILE? Close existing database and reopen FILE
                             The --new option starts with an empty file
    .output ?FILENAME?     Send output to FILENAME or stdout
    .print STRING...       Print literal STRING
    .prompt MAIN CONTINUE  Replace the standard prompts
    .quit                  Exit this program
    .read FILENAME         Execute SQL in FILENAME
    .restore ?DB? FILE     Restore content of DB (default "main") from FILE
    .save FILE             Write in-memory database into FILE
    .scanstats on|off      Turn sqlite3_stmt_scanstatus() metrics on or off
    .schema ?PATTERN?      Show the CREATE statements matching PATTERN
                              Add --indent for pretty-printing
    .selftest ?--init?     Run tests defined in the SELFTEST table
    .separator COL ?ROW?   Change the column separator and optionally the row
                             separator for both the output mode and .import
    .sha3sum ?OPTIONS...?  Compute a SHA3 hash of database content
    .shell CMD ARGS...     Run CMD ARGS... in a system shell
    .show                  Show the current values for various settings
    .stats ?on|off?        Show stats or turn stats on or off
    .system CMD ARGS...    Run CMD ARGS... in a system shell
    .tables ?TABLE?        List names of tables
                             If TABLE specified, only list tables matching
                             LIKE pattern TABLE.
    .testcase NAME         Begin redirecting output to 'testcase-out.txt'
    .timeout MS            Try opening locked tables for MS milliseconds
    .timer on|off          Turn SQL timer on or off
    .trace FILE|off        Output each SQL statement as it is run
    .vfsinfo ?AUX?         Information about the top-level VFS
    .vfslist               List all available VFSes
    .vfsname ?AUX?         Print the name of the VFS stack
    .width NUM1 NUM2 ...   Set column widths for "column" mode
                             Negative values right-justify
    sqlite> 
