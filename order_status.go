package ibapi

import "strings"

// OrderStatus represents the status of an order.
type OrderStatus int

const (
	OrderStatusApiPending    OrderStatus = iota
	OrderStatusApiCancelled
	OrderStatusPreSubmitted
	OrderStatusPendingCancel
	OrderStatusCancelled
	OrderStatusSubmitted
	OrderStatusFilled
	OrderStatusInactive
	OrderStatusPendingSubmit
	OrderStatusUnknown
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusApiPending:
		return "ApiPending"
	case OrderStatusApiCancelled:
		return "ApiCancelled"
	case OrderStatusPreSubmitted:
		return "PreSubmitted"
	case OrderStatusPendingCancel:
		return "PendingCancel"
	case OrderStatusCancelled:
		return "Cancelled"
	case OrderStatusSubmitted:
		return "Submitted"
	case OrderStatusFilled:
		return "Filled"
	case OrderStatusInactive:
		return "Inactive"
	case OrderStatusPendingSubmit:
		return "PendingSubmit"
	default:
		return "Unknown"
	}
}

// OrderStatusFromString parses a string into an OrderStatus.
func OrderStatusFromString(s string) OrderStatus {
	switch strings.ToLower(s) {
	case "apipending":
		return OrderStatusApiPending
	case "apicancelled":
		return OrderStatusApiCancelled
	case "presubmitted":
		return OrderStatusPreSubmitted
	case "pendingcancel":
		return OrderStatusPendingCancel
	case "cancelled":
		return OrderStatusCancelled
	case "submitted":
		return OrderStatusSubmitted
	case "filled":
		return OrderStatusFilled
	case "inactive":
		return OrderStatusInactive
	case "pendingsubmit":
		return OrderStatusPendingSubmit
	default:
		return OrderStatusUnknown
	}
}

// IsActive reports whether the order is in an active state.
func (s OrderStatus) IsActive() bool {
	return s == OrderStatusPreSubmitted ||
		s == OrderStatusPendingCancel ||
		s == OrderStatusSubmitted ||
		s == OrderStatusPendingSubmit
}

// IsTerminal reports whether the order has reached a terminal state.
func (s OrderStatus) IsTerminal() bool {
	return s == OrderStatusFilled ||
		s == OrderStatusCancelled ||
		s == OrderStatusInactive ||
		s == OrderStatusApiCancelled
}
