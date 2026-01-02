package enum

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusApproved PaymentStatus = "APPROVED"
	PaymentStatusRejected PaymentStatus = "REJECTED"
)

func (p PaymentStatus) String() string {
	return string(p)
}
