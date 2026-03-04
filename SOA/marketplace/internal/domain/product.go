package domain

import "time"

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Stock       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Product) IsAvailable(quantity int) bool {
	return p.Stock >= quantity
}

func (p *Product) DecreaseStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidInput
	}
	if p.Stock < quantity {
		return ErrProductOutOfStock
	}
	p.Stock -= quantity
	return nil
}

func (p *Product) IncreaseStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidInput
	}
	p.Stock += quantity
	return nil
}
