// Package mptemplates contains the HTML templates used by the application.
package mptemplates

import (
	"go/build"
	"html/template"
	"path/filepath"
)

// getSourceDir gets the directory that the source files for this package are
// installed to.
func getSourceDir() (dir string) {
	pkginfo, err := build.Import("github.com/kierdavis/mealplanner/mptemplates", "", build.FindOnly)
	if err != nil {
		panic(err)
	}

	return pkginfo.Dir
}

// The directory that the templates are stored in.
var TemplatesDir = filepath.Join(getSourceDir(), "templates")

// The parsed template object. See also: documentation on 'html/template'.
var Templates = template.Must(template.ParseGlob(filepath.Join(TemplatesDir, "*")))
