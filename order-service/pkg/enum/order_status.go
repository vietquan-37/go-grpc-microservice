package enum

type Status string

const (
	PENDING   Status = "PENDING"
	COMPLETED Status = "COMPLETED"
	CANCELLED Status = "CANCELLED"
	REFUNDED  Status = "REFUNDED"
)
