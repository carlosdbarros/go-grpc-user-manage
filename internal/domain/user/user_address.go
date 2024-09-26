package user

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Address struct {
	Street     string `json:"street"`
	Number     string `json:"number"`
	Complement string `json:"complement"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	ZipCode    string `json:"zipCode"`
}

type UserAddress struct {
	Name      string     `json:"name"`
	Emails    []string   `json:"emails"`
	Phones    []string   `json:"phones"`
	Addresses []*Address `json:"addresses"`
}

func NewAddress(street, number, complement, city, state, country, zipCode string) (*Address, error) {
	if street == "" {
		return nil, status.Error(codes.InvalidArgument, "street is required")
	}
	if number == "" {
		return nil, status.Error(codes.InvalidArgument, "number is required")
	}
	if city == "" {
		return nil, status.Error(codes.InvalidArgument, "city is required")
	}
	if state == "" {
		return nil, status.Error(codes.InvalidArgument, "state is required")
	}
	if country == "" {
		return nil, status.Error(codes.InvalidArgument, "country is required")
	}
	if zipCode == "" {
		return nil, status.Error(codes.InvalidArgument, "zipCode is required")
	}
	return &Address{
		Street:     street,
		Number:     number,
		Complement: complement,
		City:       city,
		State:      state,
		Country:    country,
		ZipCode:    zipCode,
	}, nil
}

func NewUserAddress(name string, emails []string, phones []string, addresses []*Address) (*UserAddress, error) {
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if len(emails) == 0 {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if len(phones) == 0 {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}
	if len(addresses) == 0 {
		return nil, status.Error(codes.InvalidArgument, "address is required")
	}
	var addressesResponse []*Address
	for _, a := range addresses {
		address, err := NewAddress(a.Street, a.Number, a.Complement, a.City, a.State, a.Country, a.ZipCode)
		if err != nil {
			return nil, err
		}
		addressesResponse = append(addressesResponse, address)
	}
	return &UserAddress{
		Name:      name,
		Emails:    emails,
		Phones:    phones,
		Addresses: addressesResponse,
	}, nil
}
