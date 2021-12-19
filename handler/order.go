package handler

import (
	"fmt"

	"github.com/devesh44/gin-poc/service/order"
)

type OrderHandler struct {
}

func NewOrderHandler() {
	order.New()
	fmt.Println("a")
}
