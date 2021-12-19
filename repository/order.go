package repository

type OrderRepo interface {
}

type orderRepo struct {
}

func NewOrderRepo() OrderRepo {
	return &orderRepo{}
}
