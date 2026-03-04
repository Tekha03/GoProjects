package domain

import "time"

type PromoCode struct {
	ID              string
	Code            string
	DiscountPercent int
	ExpiresAt       *time.Time
	UsageLimit      *int
	UsedCount       int
	Active          bool
	CreatedAt       time.Time
}

func (p *PromoCode) IsExpired(now time.Time) bool {
	if p.ExpiresAt == nil {
		return false
	}
	return now.After(*p.ExpiresAt)
}

func (p *PromoCode) CanBeUsed(now time.Time) error {
	if !p.Active {
		return ErrInvalidInput
	}
	if p.IsExpired(now) {
		return ErrPromoExpired
	}
	if p.UsageLimit != nil && p.UsedCount >= *p.UsageLimit {
		return ErrPromoUsageExceeded
	}
	return nil
}

func (p *PromoCode) ApplyDiscount(amount float64) float64 {
	discount := amount * float64(p.DiscountPercent) / 100
	return amount - discount
}
