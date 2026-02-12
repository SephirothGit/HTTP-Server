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

func (e OrderStatusChanged) EventName() string {
	return "order.status_changed"
}

// Entity
type Order struct {
	ID      string
	Status  string
	Version int

	events []Event
}

// Behavior
func (o *Order) ChangeStatus(newStatus string) error {
	if o.Status == newStatus {
		return nil
	}

	next, ok := ValidTransitions[o.Status]
	if !ok || !next[newStatus] {
		return nil
	}

	old := o.Status
	o.Status = newStatus
	o.Version++

	o.events = append(o.events, OrderStatusChanged{
		OrderID: o.ID,
		From:    old,
		To:      newStatus,
		At:      time.Now(),
	})

	return nil
}

func (o *Order) PullEvents() []Event {
	ev := o.events
	o.events = nil
	return ev
}
