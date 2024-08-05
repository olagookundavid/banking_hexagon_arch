package services

import (
	"github.com/olagookundavid/banking_hexagon/internal/domain"
)

type Services struct {
	Customers CustomerService
}

func NewServices() Services {
	return Services{
		Customers: NewCustomerService(domain.NewCustomerRepositoryStub()),
	}
}
