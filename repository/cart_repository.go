package repository

import (
	"database/sql"
	"fmt"
	"github.com/maurerit/shopping-cart-gomicro/proto"
	"log"
)

//TODO: Candidate for promotion
//Maybe try to mimic java's inner exceptions... will that even be useful?
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

// Finds a cart that is in a status of shopping (status = 0) and returns it.
// If no cart is found, one is created and inserted.
func GetShoppingCartByCustomerIdAndStatusOfShopping(customerId int64) (cartservice.ShoppingCart, error) {
	var rows *sql.Rows
	var cart cartservice.ShoppingCart
	var err error

	var (
		shoppingCartId int64
		status         uint32
	)
	rows, _ = DB.Query(fmt.Sprintf("select shopping_cart_id, status from shopping_cart where customer_id = %d", customerId))

	if rows.Next() {
		err = rows.Scan(&shoppingCartId, &status)
		defer rows.Close()

		if rows.Next() == true {
			err = repoError{message: "More than one row returned", innerError: err}
		}
	} else {
		insert, err := TransactionalExec("insert into shopping_cart (customer_id, status) values(?, ?)", customerId, 0)
		if err != nil {
			err = repoError{message: "Error while inserting record", innerError: err}
		}
		shoppingCartId, _ = insert.LastInsertId()
	}

	cart.CustomerId = customerId
	cart.ShoppingCartId = shoppingCartId
	cart.Status = status

	return cart, err
}

//<editor-fold desc="Prepare statement and open transaction">
//TODO: Candidate for shared function for all my services.
func TransactionalExec(sql string, args ...interface{}) (sql.Result, error) {
	var funcError error

	//TODO: I don't like the constant reassignment of funcError without first checking if funcError is nil... but I'm tired...
	transaction, err := DB.Begin()
	if err != nil {
		log.Fatal("Error opening transaction: " + err.Error())
		funcError = repoError{message: "Error opening transaction", innerError: err}
	}

	defer transaction.Rollback()
	stmt, err := transaction.Prepare(sql)
	if err != nil {
		log.Fatal("Error preparing statement: " + err.Error())
		funcError = repoError{message: "Error preparing statement", innerError: err}
	}

	result, err := stmt.Exec(args...)

	err = transaction.Commit()
	if err != nil {
		log.Fatal("Error commiting transaction")
		funcError = repoError{message: "Error committing transaction", innerError: err}
	}

	return result, funcError
}

//</editor-fold>
