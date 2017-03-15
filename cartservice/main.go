package main

import (
	"github.com/maurerit/shopping-cart-gomicro/cartservice/handling"
	"github.com/maurerit/shopping-cart-gomicro/cartservice/proto"
	"github.com/micro/go-micro"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.shoppingcart"),
	)

	service.Init()
	cartservice.RegisterCartServiceHandler(service.Server(),
		&handling.CartService{service.Client()})
	service.Run()
}
