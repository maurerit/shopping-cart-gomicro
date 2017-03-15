package handling

import (
	"fmt"
	"github.com/maurerit/shopping-cart-gomicro/cartservice/proto"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type CartService struct {
	Client client.Client
}

func (cs CartService) GetShoppingCart(ctx context.Context,
	request *cartservice.ShoppingCartRequest,
	response *cartservice.ShoppingCart) error {

	fmt.Println(request)
	response.ShoppingCartId = 2000
	response.CustomerId = 2000
	response.Status = 2000
	return nil
}

func (cs CartService) GetByCartId(ctx context.Context,
	request *cartservice.ShoppingCartRequest,
	response *cartservice.ShoppingCart) error {

	return nil
}

func (cs CartService) AddItemToCart(ctx context.Context,
	request *cartservice.ShoppingCartItem,
	response *cartservice.ApiResponse) error {

	return nil
}

func (cs CartService) Checkout(ctx context.Context,
	request *cartservice.ShoppingCartRequest,
	response *cartservice.ApiResponse) error {

	return nil
}

func (cs CartService) UpdateItem(ctx context.Context,
	request *cartservice.ShoppingCartItem,
	response *cartservice.ApiResponse) error {

	return nil
}
