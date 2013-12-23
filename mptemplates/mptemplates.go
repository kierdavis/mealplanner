package mptemplates

import (
    "go/build"
    "html/template"
    "path/filepath"
)

func getSourceDir() (dir string) {
    pkginfo, err := build.Import("github.com/kierdavis/mealplanner/mptemplates", "", build.FindOnly)
    if err != nil {
    	panic(err)
    }
    
    return pkginfo.Dir
}

var TemplatesDir = filepath.Join(getSourceDir(), "templates")

var Templates = template.Must(template.ParseGlob(filepath.Join(TemplatesDir, "*")))
