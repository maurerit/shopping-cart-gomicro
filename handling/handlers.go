package handling

import (
	"github.com/maurerit/shopping-cart-gomicro/proto"
	"github.com/maurerit/shopping-cart-gomicro/repository"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type CartService struct {
	Client client.Client
}

type handlerError struct {
	message string
}

func (he handlerError) Error() string {
	return he.message
}

func (cs CartService) GetShoppingCart(ctx context.Context,
	request *cartservice.ShoppingCartRequest,
	response *cartservice.ShoppingCart) error {

	cart, err := repository.GetShoppingCartByCustomerIdAndStatusOfShopping(request.CustomerId)
	if err != nil {
		cart, _ = repository.CreateShoppingCartForCustomer(request.CustomerId)
	}

	cartItems, err := repository.GetShoppingCartItemsByShoppingCartId(cart.ShoppingCartId)

	response.ShoppingCartId = cart.ShoppingCartId
	response.CustomerId = cart.CustomerId
	response.Items = cartItems
	response.Status = cart.Status

	return nil
}

func (cs CartService) GetByCartId(ctx context.Context,
	request *cartservice.ShoppingCartRequest,
	response *cartservice.ShoppingCart) error {

	cart, err := repository.GetShoppingCartByShoppingCartId(request.ShoppingCartId)

	if err == nil {
		response.ShoppingCartId = cart.ShoppingCartId
		response.CustomerId = cart.CustomerId
		response.Items, _ = repository.GetShoppingCartItemsByShoppingCartId(cart.ShoppingCartId)
		response.Status = cart.Status
	} else {
		return handlerError{message: "GetByCartId produced an error: " + err.Error()}
	}

	return nil
}

func (cs CartService) AddItemToCart(ctx context.Context,
	request *cartservice.ShoppingCartItem,
	response *cartservice.ApiResponse) error {

	_, err := repository.GetShoppingCartByShoppingCartId(request.ShoppingCartId)

	if err == nil {
		err = repository.CreateShoppingCartItemForShopingCart(request)

		if err == nil {
			response.Success = "true"
		} else {
			response.Success = "false"
		}
	} else {
		response.Success = "false"
	}

	return err
}

func (cs CartService) Checkout(ctx context.Context,
	request *cartservice.ShoppingCartRequest,
	response *cartservice.ApiResponse) error {

	cart, err := repository.GetShoppingCartByShoppingCartId(request.ShoppingCartId)
	if err == nil {
		cartItems, err := repository.GetShoppingCartItemsByShoppingCartId(request.ShoppingCartId)
		if err == nil {
			cart.Status = repository.ORDERED
			err = repository.UpdateShoppingCart(cart)

			for _, cartItem := range cartItems {
				cartItem.Status = repository.ORDERED
				repository.UpdateShoppingCartItem(cartItem)
			}

			response.Success = "true"
		} else {
			response.Success = "false"
		}
	} else {
		response.Success = "false"
	}

	return err
}

func (cs CartService) UpdateItem(ctx context.Context,
	request *cartservice.ShoppingCartItem,
	response *cartservice.ApiResponse) error {

	return nil
}
