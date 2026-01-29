package domain

import "errors"

//Business errors
var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidTransition = errors.New("invalid status transition")
)

//Domain entity
type Order struct {
	ID string
	Status string
}

//Allowed statuses
const (
	StatusCreated = "created"
	StatusPaid = "paid"
	StatusShipped = "shipped"
	StatusCanceled = "canceled"
)

//Valid transitions that describes allowed status changes
var ValidTransitions = map[string]map[string]bool{
	StatusCreated: {
		StatusPaid: true,
		StatusCanceled: true,
	},
	StatusPaid: {
		StatusShipped: true,
		StatusCanceled: true,
	},
	StatusShipped: {},
	StatusCanceled: {},
}

//Cantrsnsition checks if status change is allowed
func CanTransition(from, to string) bool {
	next, ok := ValidTransitions[from]
	if !ok {
		return false
	}
	return next[to]
}