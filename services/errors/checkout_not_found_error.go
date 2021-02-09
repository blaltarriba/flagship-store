package errors

type CheckoutNotFoundError struct {
	data string
}

func NewCheckoutNotFoundError() error {
	return &CheckoutNotFoundError{}
}

func (e *CheckoutNotFoundError) Error() string {
	return ""
}
