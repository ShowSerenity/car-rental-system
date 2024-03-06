package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	mux := pat.New()

	// Home page
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	// Simplified car routes using a loop
	carTypes := []string{"economy", "comfort", "minivan"}
	for _, t := range carTypes {
		mux.Get("/cars/"+t, dynamicMiddleware.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
			app.showCarsByType(w, r, t)
		}))
	}

	// Car creation and display routes
	mux.Get("/cars/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createCarForm))
	mux.Post("/cars/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createCar))
	mux.Get("/cars/:id", dynamicMiddleware.ThenFunc(app.showCar))

	// Rent routes
	mux.Get("/rent/create/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createRentForm))
	mux.Post("/rent/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createRent))
	mux.Get("/rent/:id", dynamicMiddleware.ThenFunc(app.showRent))

	// User routes
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))
	mux.Get("/user/profile/", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.profileUser))

	// Static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
