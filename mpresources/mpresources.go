// Package mpresources contains the HTML templates and static files used by the
// application.
package mpresources

import (
	"go/build"
	"html/template"
	"path/filepath"
)

// resourceDir is the directory that all resources are stored in. It should be
// considered uninitialised until GetResourceDir is called or it is otherwise
// assigned to.
var resourceDir string

// getSourceDir gets the directory that the source files for this package are
// installed to.
func getSourceDir() (dir string) {
	pkginfo, err := build.Import("github.com/kierdavis/mealplanner/mpresources", "", build.FindOnly)
	if err != nil {
		panic("Resource directory not set and no suitable directory found in the GOPATH")
	}

	return pkginfo.Dir
}

// GetResourceDir returns the resource directory. If it is uninitialised, it
// looks for the package's source directory in the GOPATH and uses that.
func GetResourceDir() (dir string) {
	if resourceDir == "" {
		resourceDir = filepath.Join(getSourceDir(), "resources")
	}

	return resourceDir
}

// SetResourceDir sets the resource directory.
func SetResourceDir(dir string) {
	resourceDir = dir
}

// GetStaticDir returns the directory used for storing static files.
func GetStaticDir() (dir string) {
	return filepath.Join(GetResourceDir(), "static")
}

var templates *template.Template

// GetTemplates loads the templates from the resource directory if they have not
// been loaded already, and returns them.
func GetTemplates() (t *template.Template) {
	if templates == nil {
		pattern := filepath.Join(GetResourceDir(), "templates", "*")
		templates = template.Must(template.ParseGlob(pattern))
	}

	return templates
}

/*
// ResourcesDir is the directory that all resources are stored in.
var ResourcesDir = filepath.Join(getSourceDir(), "resources")

// TemplatesDir is the directory that the templates are stored in.
var TemplatesDir = filepath.Join(ResourcesDir, "templates")

// StaticDir is the directory that static files are stored in.
var StaticDir = filepath.Join(ResourcesDir, "static")

// Templates contains the parsed templates. See also: documentation on
// 'html/template'.
var Templates = template.Must(template.ParseGlob(filepath.Join(TemplatesDir, "*")))
*/
