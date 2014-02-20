Meal Planner (A2 Computing Project)

TO-DO: write a description here

# Installation & Usage

Download & install with:

    go get github.com/kierdavis/mealplanner

Run server with:

    mealplanner -dbsource 'DBSOURCE'

where `DBSOURCE` is a data source identifier of the form:

    USER:PASS@tcp(HOST:PORT)/DATABASE
    or
    USER:PASS@unix(/PATH/TO/SOCKET)/DATABASE

For example, to connect the local 'mealplanner' database with username
'web' and password 'hello123', the data source would be:

    web:hello123@unix(/var/run/mysqld/mysqld.sock)/mealplanner

and to connect with the same credentials to a remote database at db.example.net,
port 3306 (the default MySQL port):

    web:hello123@tcp(db.example.net:3306)/mealplanner

Alternatively to using the `-dbsource` flag, the `MPDBSOURCE` environment
variable can also be set. The `-dbsource` flag overrides the environment
variable. Example:

    export MPDBSOURCE='web:hello123@tcp(db.example.net:3306)/mealplanner'
    mealplanner

# API documentation

Hosted on [GoDoc](http://godoc.org/github.com/kierdavis/mealplanner).

# Code structure

* `github.com/kierdavis/mealplanner` - Main application command
* `github.com/kierdavis/mealplanner/mpapi` - Implementations of Ajax API calls
* `github.com/kierdavis/mealplanner/mpdata` - Data structures and data-
  processing algorithms
* `github.com/kierdavis/mealplanner/mpdb` - Database interface (abstraction of
  SQL commands into functions involving the application's data structures)
* `github.com/kierdavis/mealplanner/mphandlers` - Implementations of HTTP
  request handlers
* `github.com/kierdavis/mealplanner/mpresources` - Data files (HTML templates,
  static web files), and code to locate these at runtime
