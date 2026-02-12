package domain

import (
	"errors"
	"time"
)

// Errors
var (
	ErrOrderNotFound     = errors.New("Order not found")
	ErrInvalidTransition = errors.New("Invalid status transition")
	ErrVersionConflict   = errors.New("Version conflict")
)
// Domain entity
type Order struct {
	ID      string
	Status  string
	Version int
	events []Event
}

// Statuses
const (
	StatusCreated  = "created"
	StatusPaid     = "paid"
	StatusShipped  = "shipped"
	StatusCanceled = "canceled"
)

// Transitions
var ValidTransitions = map[string]map[string]bool{
	StatusCreated: {
		StatusPaid:     true,
		StatusCanceled: true,
	},
	StatusPaid: {
		StatusShipped:  true,
		StatusCanceled: true,
	},
	StatusShipped:  {},
	StatusCanceled: {},
}

// Domain event
type Event interface {
	EventName() string
}
type OrderStatusChanged struct {
	OrderID string
	From    string
	To      string
	At      time.Time
}
// Checks if status change allowed
func CanTransition(from, to string) bool {
	next, ok := ValidTransitions[from]
	if !ok {
		return false
	}
	return next[to]
}

// Checks if the order already have the given status
func IsSameStatus(current, new string) bool {
	return current == new
}


func (o *Order) ChangeStatus(newStatus string) error {

	// Idempotency check
	if IsSameStatus(o.Status, newStatus) {
		return nil
	}

	// Valid transition check
	if !CanTransition(o.Status, newStatus) {
		return ErrInvalidTransition
	}

	// Update status
	from := o.Status
	o.Status = newStatus
	o.Version++

	// Record a domain event
	o.events = append(o.events, OrderStatusChanged{
		OrderID: o.ID,
		From: from,
		To: newStatus,
		At: time.Now(),
	})

	return nil
}

func (e OrderStatusChanged) EventName() string {
	return "order.status_changed"
}

func (o *Order) PullEvents() []Event {
	ev := o.events
	o.events = nil
	return ev
}