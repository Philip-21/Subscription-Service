package main

import (
	"fmt"
	"html/template"
	"net/http"
	"subscription-service/database"
	"time"
)

const pathToTemplates = "./cmd/web/templates"

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	Data          map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	User          *database.User
}

func (app *Config) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {
	partials := []string{
		fmt.Sprintf("%s/base.layout.go.html", pathToTemplates),
		fmt.Sprintf("%s/header.partial.go.html", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.go.html", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.go.html", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.go.html", pathToTemplates),
	}

	var templateSlice []string
	//renders the main template
	templateSlice = append(templateSlice, fmt.Sprintf("%s/%s", pathToTemplates, t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	if td == nil {
		td = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, app.AddDefaultData(td, r)); err != nil {
		app.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// describes the data in the template in reference to templatedata
func (app *Config) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	if app.IsAuthenticated(r) {
		td.Authenticated = true
		// get more user information if authenticated and put in a session
		user, ok := app.Session.Get(r.Context(), "user").(database.User)
		if !ok {
			app.ErrorLog.Println("Cant get User from Session")
		} else {
			//add the user to template info
			td.User = &user
		}
	}
	//current date and time
	td.Now = time.Now()

	return td
}
