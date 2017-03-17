package repository

import (
	"database/sql"
	"fmt"
	"github.com/maurerit/shopping-cart-gomicro/proto"
	"strconv"
)

// Finds all shopping cart items that are in the shopping cart identified by the id.
// Returns all the items and an errors that happened during the query or nil if nothing went wrong.
func GetShoppingCartItemsByShoppingCartId(shoppingCartId int64) ([]*cartservice.ShoppingCartItem, error) {
	var rows *sql.Rows
	var cartItems []*cartservice.ShoppingCartItem
	var err error

	cartItems = make([]*cartservice.ShoppingCartItem, 0, 1000)

	var (
		itemId   int64
		quantity int64
		price    float32
		status   uint32
	)
	rows, err = DB.Query(fmt.Sprintf("select item_id, quantity, price, status from shopping_cart_item where shopping_cart_id = %d ", shoppingCartId))

	defer rows.Close()
	for idx := 0; rows.Next(); idx++ {
		if len(cartItems) >= idx {
			rows.Scan(&itemId, &quantity, &price, &status)
			cartItems = append(cartItems, &cartservice.ShoppingCartItem{
				ShoppingCartId: shoppingCartId,
				ItemId:         itemId,
				Quantity:       quantity,
				Price:          price,
				Status:         status,
			})
		}
	}

	return cartItems, err
}

// Creates a new shopping cart item and returns any errors if something went wrong or nil
func CreateShoppingCartItemForShopingCart(cartItem *cartservice.ShoppingCartItem) error {

	_, err := TransactionalExec(
		"insert into shopping_cart_item (shopping_cart_id, item_id, quantity, price, status) values(?, ?, ?, ?, ?)",
		cartItem.ShoppingCartId,
		cartItem.ItemId,
		cartItem.Quantity,
		cartItem.Price,
		SHOPPING)

	if err != nil {
		err = repoError{message: "Error while inserting record", innerError: err}
	}

	return err
}

// Updates the given shopping cart item.
func UpdateShoppingCartItem(cartItem *cartservice.ShoppingCartItem) error {
	result, err := TransactionalExec(
		"update shopping_cart_item set item_id = ?, quantity = ?, price = ?, status = ? where shopping_cart_id = ? and item_id = ?",
		cartItem.ItemId,
		cartItem.Quantity,
		cartItem.Price,
		cartItem.Status,
		cartItem.ShoppingCartId,
		cartItem.ItemId)

	//This should really never happen as shopping_cart_id and item_id are the pkey... but who knows...
	if err == nil {
		rows, err := result.RowsAffected()
		err = repoError{message: "Multiple rows affected: " + strconv.FormatInt(rows, 10), innerError: err}
	}

	return err
}
