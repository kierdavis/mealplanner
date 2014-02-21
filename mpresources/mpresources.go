// Package mpresources contains the HTML templates and static files used by the
// application.
package mpresources

import (
	"go/build"
	"html/template"
	"path/filepath"
)

// getSourceDir gets the directory that the source files for this package are
// installed to.
func getSourceDir() (dir string) {
	pkginfo, err := build.Import("github.com/kierdavis/mealplanner/mpresources", "", build.FindOnly)
	if err != nil {
		panic(err)
	}

	return pkginfo.Dir
}

// ResourcesDir is the directory that all resources are stored in.
var ResourcesDir = filepath.Join(getSourceDir(), "resources")

// TemplatesDir is the directory that the templates are stored in.
var TemplatesDir = filepath.Join(ResourcesDir, "templates")

// StaticDir is the directory that static files are stored in.
var StaticDir = filepath.Join(ResourcesDir, "static")

// Templates contains the parsed templates. See also: documentation on
// 'html/template'.
var Templates = template.Must(template.ParseGlob(filepath.Join(TemplatesDir, "*")))
