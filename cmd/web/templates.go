package main

import (
	"html/template"
	"path/filepath"
	"showserenity.net/car-rental-system/pkg/forms"
	"showserenity.net/car-rental-system/pkg/models"
	"time"
)

type templateData struct {
	CSRFToken       string
	CurrentYear     int
	Flash           string
	Form            *forms.Form
	Car             *models.Car
	Cars            []*models.Car
	Rent            *models.Rent
	Rents           []*models.Rent
	User            *models.User
	Users           []*models.User
	IsAuthenticated bool
	IsAdmin         bool
	Error           string
	CarsType        string
	IframeSrc       string
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	layoutFiles, err := filepath.Glob(filepath.Join(dir, "*.layout.tmpl"))
	if err != nil {
		return nil, err
	}

	partialFiles, err := filepath.Glob(filepath.Join(dir, "*.partial.html"))
	if err != nil {
		return nil, err
	}

	// combine layout and partial files which will be included in every page template
	templateFiles := append(layoutFiles, partialFiles...)

	pageFiles, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pageFiles {
		name := filepath.Base(page)

		// parse the page template along with the layout and partial templates
		ts, err := template.New(name).Funcs(functions).ParseFiles(append(templateFiles, page)...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
