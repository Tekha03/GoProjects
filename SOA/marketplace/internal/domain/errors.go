package domain

import "net/http"

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrInternal = &AppError{
		Code:       "internal_error",
		Message:    "internal server error",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrNotFound = &AppError{
		Code:       "not_found",
		Message:    "resource not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrUnauthorized = &AppError{
		Code:       "unauthorized",
		Message:    "unauthorized",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrForbidden = &AppError{
		Code:       "forbidden",
		Message:    "forbidden",
		HTTPStatus: http.StatusForbidden,
	}

	ErrInvalidInput = &AppError{
		Code:       "invalid_input",
		Message:    "invalid input",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrProductOutOfStock = &AppError{
		Code:       "product_out_of_stock",
		Message:    "product out of stock",
		HTTPStatus: http.StatusConflict,
	}

	ErrInvalidOrderStatusTransition = &AppError{
		Code:       "invalid_order_status_transition",
		Message:    "invalid order status transition",
		HTTPStatus: http.StatusConflict,
	}

	ErrPromoExpired = &AppError{
		Code:       "promo_expired",
		Message:    "promo code expired",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrPromoUsageExceeded = &AppError{
		Code:       "promo_usage_exceeded",
		Message:    "promo usage limit exceeded",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrEmailAlreadyExists = &AppError{
		Code:       "email_already_exists",
		Message:    "email already exists",
		HTTPStatus: http.StatusConflict,
	}
)
