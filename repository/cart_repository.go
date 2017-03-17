package repository

import (
	"database/sql"
	"fmt"
	"github.com/maurerit/shopping-cart-gomicro/proto"
	"strconv"
)

// Finds a cart that is in a status of shopping (status = 1) and returns it.
// If no cart is found, an error is returned.
func GetShoppingCartByCustomerIdAndStatusOfShopping(customerId int64) (*cartservice.ShoppingCart, error) {
	var rows *sql.Rows
	var cart *cartservice.ShoppingCart
	var err error

	var (
		shoppingCartId int64
		status         uint32
	)
	rows, _ = DB.Query(fmt.Sprintf("select shopping_cart_id, status from shopping_cart where customer_id = %d and status = %d", customerId, SHOPPING))

	if rows.Next() {
		err = rows.Scan(&shoppingCartId, &status)
		defer rows.Close()

		if rows.Next() == true {
			err = repoError{message: "More than one row returned", innerError: err}
		}

		cart = &cartservice.ShoppingCart{CustomerId: customerId, ShoppingCartId: shoppingCartId, Status: status}
	} else {
		return nil, repoError{message: "No cart in status shopping found for that customer", innerError: nil}
	}

	return cart, err
}

// Fetches the shopping cart using the given unique id.  If none exist then a nil cart and error are returned.
func GetShoppingCartByShoppingCartId(shoppingCartId int64) (*cartservice.ShoppingCart, error) {
	var rows *sql.Rows
	var cart *cartservice.ShoppingCart
	var err error

	var (
		customerId int64
		status     uint32
	)
	rows, _ = DB.Query(fmt.Sprintf("select customer_id, status from shopping_cart where shopping_cart_id = %d", shoppingCartId))

	if rows.Next() {
		err = rows.Scan(&customerId, &status)
		defer rows.Close()

		if rows.Next() == true {
			err = repoError{message: "More than one row returned", innerError: err}
		}

		cart = &cartservice.ShoppingCart{CustomerId: customerId, ShoppingCartId: shoppingCartId, Status: status}
	} else {
		return nil, repoError{message: "Could not find a cart by that id", innerError: nil}
	}

	return cart, err
}

// Creates a shopping cart id in the status of shopping.  If there was an error inserting or fetching the auto_incremented
// shopping_cart_id then a nil cart and error are returned.
func CreateShoppingCartForCustomer(customerId int64) (*cartservice.ShoppingCart, error) {
	var cart *cartservice.ShoppingCart

	insert, err := TransactionalExec("insert into shopping_cart (customer_id, status) values(?, ?)", customerId, SHOPPING)
	if err == nil {
		shoppingCartId, err := insert.LastInsertId()

		if err == nil {
			cart = &cartservice.ShoppingCart{CustomerId: customerId, ShoppingCartId: shoppingCartId, Status: SHOPPING}
		} else {
			err = repoError{message: "Error grabbing last inserted id", innerError: err}
		}
	} else {
		return nil, repoError{message: "Error while inserting record", innerError: err}
	}

	return cart, err
}

// Updates a shopping cart setting all the fields to the fields in the passed in struct.  If there is an error then it's returned.
func UpdateShoppingCart(cart *cartservice.ShoppingCart) error {
	result, err := TransactionalExec(
		"update shopping_cart set customer_id = ?, status = ? where shopping_cart_id = ?",
		cart.CustomerId,
		cart.Status,
		cart.ShoppingCartId)

	//TODO: I don't really like this... make it prettier
	//This should really never happen as shopping_cart_id is the pkey... but who knows...
	if err == nil {
		rows, err := result.RowsAffected()
		//Use err as it could be nil fromt he RowsAffected function call above
		return repoError{message: "Multiple rows affected: " + strconv.FormatInt(rows, 10), innerError: err}
	} else {
		return repoError{message: "Updating cart failed", innerError: err}
	}
}
