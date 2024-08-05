package app

import "net/http"

type CustomerHandlers struct {
	// service services.CustomerService
}

func (app *Application) GetAllCustomers(w http.ResponseWriter, r *http.Request) {

	customers, _ := app.Services.Customers.GetAllCustomer("")
	env := envelope{
		"message":   "Retrieved All Customers",
		"customers": customers}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) GetCustomerById(w http.ResponseWriter, r *http.Request) {

	customer, _ := app.Services.Customers.GetCustomer("")
	env := envelope{
		"message":  "Retrieved All Customer",
		"customer": customer}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
