package main

import (
	"errors"
	"fmt"
	"net/http"
	"showserenity.net/car-rental-system/pkg/forms"
	"showserenity.net/car-rental-system/pkg/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	/*if r.URL.Path != "/" {
		app.notFound(w)
		return
	}*/

	c, err := app.cars.LatestCars()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.html", &templateData{Cars: c})

}

func (app *application) showEconomy(w http.ResponseWriter, r *http.Request) {
	app.showCarsByType(w, r, "economy")
}
func (app *application) showComfort(w http.ResponseWriter, r *http.Request) {
	app.showCarsByType(w, r, "comfort")
}
func (app *application) showMinivan(w http.ResponseWriter, r *http.Request) {
	app.showCarsByType(w, r, "minivan")
}

func (app *application) showCarsByType(w http.ResponseWriter, r *http.Request, carsType string) {

	carsList, err := app.cars.GetCarsByType(carsType)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "cars.page.html", &templateData{CarsType: carsType, Cars: carsList})
}

func (app *application) showRent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.rents.GetRent(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	iframeSrc := app.generateRandomMap()
	s.Location = iframeSrc

	app.render(w, r, "rentShow.page.html", &templateData{Rent: s, IframeSrc: iframeSrc})
}

func (app *application) createRentForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	c, err := app.cars.GetCar(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	check, err := app.rents.GetRentByCarID(c.ID)
	if check {
		errorMessage := "This vehicle is currently rented by someone else"
		app.render(w, r, "error.page.html", &templateData{Error: errorMessage})
	}

	app.render(w, r, "rentCreate.page.html", &templateData{
		Form: forms.NewRent(nil, c)})
}

func (app *application) createRent(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.NewRent(r.PostForm, nil)

	if !form.Valid() {
		app.render(w, r, "rentCreate.page.html", &templateData{Form: form})
		return
	}

	hours, err := strconv.Atoi(form.Get("hours"))
	minutes, err := strconv.Atoi(form.Get("minutes"))
	time := (hours * 60) + minutes

	// Retrieve the carID from the form data
	carID, err := strconv.Atoi(form.Get("carID"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Retrieve the totalCost from the form data
	totalCost, err := strconv.ParseFloat(form.Get("totalCost"), 64)
	if err != nil {
		app.serverError(w, err)
		return
	}

	id, err := app.rents.InsertRent(app.session.Get(r, "authenticatedUserID").(int), carID, time, totalCost)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Rent record successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/rent/%d", id), http.StatusSeeOther)
}

func (app *application) showCar(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	c, err := app.cars.GetCar(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	iframeSrc := app.generateRandomMap()
	c.Location = iframeSrc

	app.render(w, r, "carShow.page.html", &templateData{Car: c, IframeSrc: iframeSrc})
}

func (app *application) createCarForm(w http.ResponseWriter, r *http.Request) {
	id := app.session.Get(r, "authenticatedUserID").(int)

	role, err := app.users.RoleCheck(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Allow access if the user is an admin or teacher
	if role != "admin" {
		errorMessage := "You are not authorized to access this page"
		app.render(w, r, "error.page.html", &templateData{Error: errorMessage})
		return
	}

	app.render(w, r, "carCreate.page.html", &templateData{
		Form: forms.NewCar(nil)})
}

func (app *application) createCar(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.NewCar(r.PostForm)
	form.Required("seats", "age", "cost", "model", "carType", "color", "location", "imageUrl", "description")
	form.MaxLength("model", 100)
	form.PermittedValues("seats", "2", "5", "8")
	form.PermittedValues("carType", "economy", "comfort", "minivan")

	if !form.Valid() {
		app.render(w, r, "carCreate.page.html", &templateData{Form: form})
		return
	}

	seatsStr := form.Get("seats")
	seats, err := strconv.Atoi(seatsStr)
	ageStr := form.Get("age")
	age, err := strconv.Atoi(ageStr)
	costStr := form.Get("cost")
	cost, err := strconv.Atoi(costStr)

	id, err := app.cars.InsertCar(seats, age, cost, form.Get("model"),
		form.Get("carType"), form.Get("color"),
		form.Get("location"), form.Get("imageUrl"), form.Get("description"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Rent record successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/cars/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.html", &templateData{
		Form: forms.NewSignUp(nil),
	})
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.NewSignUp(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	if !form.Valid() {
		app.render(w, r, "signup.page.html", &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.html", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.html", &templateData{
		Form: forms.NewSignUp(nil),
	})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.NewSignUp(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.html", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "authenticatedUserID", id)
	userID := app.session.Get(r, "authenticatedUserID").(int)

	role, err := app.users.RoleCheck(userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	if role == "admin" {
		app.session.Put(r, "isAdmin", true)
	} else {
		app.session.Put(r, "isAdmin", false)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Remove(r, "isAdmin")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) profileUser(w http.ResponseWriter, r *http.Request) {
	id := app.session.Get(r, "authenticatedUserID").(int)

	u, err := app.users.GetUser(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	rents, err := app.rents.LatestRents(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "profile.page.html", &templateData{
		User:  u,
		Rents: rents,
	})
}
