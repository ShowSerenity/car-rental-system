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
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	mux.Get("/rent/create/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createRentForm))
	mux.Post("/rent/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createRent))
	mux.Get("/rent/:id", dynamicMiddleware.ThenFunc(app.showRent))

	mux.Get("/cars/sedan", dynamicMiddleware.ThenFunc(app.showSedan))
	mux.Get("/cars/pickup", dynamicMiddleware.ThenFunc(app.showPickup))
	mux.Get("/cars/sportCar", dynamicMiddleware.ThenFunc(app.showSportCar))
	mux.Get("/cars/minivan", dynamicMiddleware.ThenFunc(app.showMinivan))
	mux.Get("/cars/suv", dynamicMiddleware.ThenFunc(app.showSUV))

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
