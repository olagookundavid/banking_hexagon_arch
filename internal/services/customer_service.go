package services

import (
	"github.com/olagookundavid/banking_hexagon/internal/domain"
)

type CustomerService interface {
	GetAllCustomer(status string) ([]domain.Customer, error)
	GetCustomer(string) (domain.Customer, error)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository}
}

func (s DefaultCustomerService) GetAllCustomer(status string) ([]domain.Customer, error) {
	// if status == "active" {
	// 	status = "1"
	// } else if status == "inactive" {
	// 	status = "0"
	// } else {
	// 	status = ""
	// }
	// customers, err := s.repo.FindAll(status)
	// if err != nil {
	// 	return nil, err
	// }
	// response := make([]dto.CustomerResponse, 0)
	// for _, c := range customers {
	// 	response = append(response, c.ToDto())
	// }
	return s.repo.FindAll()
}

func (s DefaultCustomerService) GetCustomer(status string) (domain.Customer, error) {

	return s.repo.ById("")
}
