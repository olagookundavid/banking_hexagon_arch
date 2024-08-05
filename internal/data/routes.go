package data

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	app "github.com/olagookundavid/banking_hexagon/cmd/api"
)

func Routes(app *app.Application) http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/customers", app.GetAllCustomers)
	router.HandlerFunc(http.MethodGet, "/v1/customer", app.GetCustomerById)
	return router
}

// type Models struct {
// 	// Ch api.CustomerHandlers
// }

// func NewModels() Models {
// 	return Models{
// 		// Ch: api.CustomerHandlers{services.NewCustomerService(domain.NewCustomerRepositoryStub())},
// 	}
// }
