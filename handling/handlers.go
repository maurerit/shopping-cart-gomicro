package handling

import (
	"database/sql"
	"github.com/maurerit/shopping-cart-gomicro/proto"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type CartService struct {
	Client client.Client
}

func (cs CartService) GetShoppingCart(ctx context.Context,
	request *cartservice.ShoppingCartRequest,
	response *cartservice.ShoppingCart) error {

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

func buildShoppingCartFromRows(rows *sql.Rows) (cartservice.ShoppingCart, error) {
	return cartservice.ShoppingCart{}, nil
}

func buildShoppingCartItemFromRows(rows *sql.Rows) (cartservice.ShoppingCartItem, error) {
	return cartservice.ShoppingCartItem{}, nil
}
