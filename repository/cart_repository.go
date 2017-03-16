package repository

import (
	"database/sql"
	"fmt"
	"github.com/maurerit/shopping-cart-gomicro/proto"
)

type repoError struct {
	message    string
	innerError error
}

func (re repoError) Error() string {
	message := re.message

	if re.innerError != nil {
		message += " caused by " + re.innerError.Error()
	}

	return message
}

// Global database handle
var DB *sql.DB

func GetShoppingCartByCustomerIdAndStatusOfShopping(customerId int64) (cartservice.ShoppingCart, error) {
	var rows *sql.Rows
	var cart cartservice.ShoppingCart

	var (
		shoppingCartId int64
		status         uint32
	)
	rows, _ = DB.Query(fmt.Sprintf("select shopping_cart_id, status from shopping_cart where customer_id = %d", customerId))

	rows.Next()
	err := rows.Scan(&shoppingCartId, &status)
	defer rows.Close()

	if rows.Next() == true {
		err = repoError{message: "More than one row returned", innerError: err}
	}

	cart.CustomerId = customerId
	cart.ShoppingCartId = shoppingCartId
	cart.Status = status

	return cart, err
}

//<editor-fold desc="Prepare statement and open transaction">
//TODO: Candidate for shared function for all my services.
func TransactionalExec(string string, args ...interface{}) (*sql.Result, []error) {
	return nil, nil
}

//</editor-fold>
