package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

type Order struct {
	ID          string
	UserID      string
	Status      OrderStatus
	Items       []OrderItem
	TotalAmount float64
	PromoID     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (o *Order) CalculateTotal() {
	total := 0.0
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	o.TotalAmount = total
}

func (o *Order) ApplyPromo(promo *PromoCode, now time.Time) error {
	if promo == nil {
		return nil
	}
	if err := promo.CanBeUsed(now); err != nil {
		return err
	}
	o.TotalAmount = promo.ApplyDiscount(o.TotalAmount)
	return nil
}

func (o *Order) ChangeStatus(newStatus OrderStatus) error {
	switch o.Status {
	case OrderStatusPending:
		if newStatus == OrderStatusPaid || newStatus == OrderStatusCancelled {
			o.Status = newStatus
			return nil
		}
	case OrderStatusPaid:
		// paid -> no transitions allowed
	case OrderStatusCancelled:
		// cancelled -> no transitions allowed
	}
	return ErrInvalidOrderStatusTransition
}
