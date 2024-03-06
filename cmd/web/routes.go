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
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	mux.Get("/cars/economy", dynamicMiddleware.ThenFunc(app.showEconomy))
	mux.Get("/cars/comfort", dynamicMiddleware.ThenFunc(app.showComfort))
	mux.Get("/cars/minivan", dynamicMiddleware.ThenFunc(app.showMinivan))

	mux.Get("/cars/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createCarForm))
	mux.Post("/cars/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createCar))
	mux.Get("/cars/:id", dynamicMiddleware.ThenFunc(app.showCar))

	mux.Get("/rent/create/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createRentForm))
	mux.Post("/rent/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createRent))
	mux.Get("/rent/:id", dynamicMiddleware.ThenFunc(app.showRent))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))
	mux.Get("/user/profile/", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.profileUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
