Meal Planner (A2 Computing Project)

TO-DO: write a description here

# Installation

Download & install with:

    go get github.com/kierdavis/mealplanner

Run server with:

    mealplanner

# Code structure

* `github.com/kierdavis/mealplanner` - Main application command
* `github.com/kierdavis/mealplanner/mpapi` - Implementations of Ajax API calls
    * `mpapi.go` - Dispatching code
    * other files - Implementations of individial API calls
* `github.com/kierdavis/mealplanner/mpdata` - Data structures and data-processing algorithms
    * `mpdata.go` - Miscellaneous: global constants
    * `types.go` - Data structure definitions
    * `score.go` - Scoring algorithm
* `github.com/kierdavis/mealplanner/mpdb` - Database interface (abstraction of SQL commands into functions 
involving the application's data structures)
    * `mpdb.go` - Connection details, definition of `Queryable` interface and closure functions (`With*`)
    * `tables.go` - Database routines to create empty database & tables
    * other files - Database routines for manipulating different areas of the database
* `github.com/kierdavis/mealplanner/mphandlers` - Implementations of HTTP request handlers
    * `mphandlers.go` - Dispatching code
    * `httperror.go` - Definition of HTTP error codes & associated user messages (this might be temporary)
    * `util.go` - Utility functions
    * other files - Implementations of individual web page handlers
* `github.com/kierdavis/mealplanner/mpresources` - Data files (HTML templates, static web files), and code to 
locate these at runtime
    * `mpresources.go` - Resource loader
    * `resources/templates/` - HTML templates
    * `resources/static/` - Static files
        * `resources/static/css/screen.css` - Application-specific CSS
        * `resources/static/js/mpajax.js` - Client-side interface to Ajax API
