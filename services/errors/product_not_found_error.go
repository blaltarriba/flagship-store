package errors

type ProductNotFoundError struct {
	data string
}

func NewProductNotFoundError() error {
	return &ProductNotFoundError{}
}

func (e *ProductNotFoundError) Error() string {
	return ""
}
