package repository

import (
	"database/sql"
	//"log"
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

func GetShoppingCartByCustomerId(customerId int64) (cartservice.ShoppingCart, error /*[5]error*/) {
	var rows *sql.Rows
	//var allErrors [5]error
	var cart cartservice.ShoppingCart

	var (
		shoppingCartId int64
		status         uint32
	)

	//rows, errors := TransactionalQuery("select * from shopping_cart where customer_id = ?", customerId)

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

	//for idx, curErr := range errors {
	//	allErrors[idx] = curErr
	//}
	//allErrors[4] = err

	return cart, err
}

//<editor-fold desc="Prepare statement and open transaction">
//TODO: Candidate for shared function for all my services.
//TODO: This isn't working no matter how much I tweak it :(
// Prepares a transaction and a statement and opens them.  Returns any errors it encountered
func TransactionalQuery(sql string, args ...interface{}) (*sql.Rows, [4]error) {
	var errors [4]error

	transaction, err := DB.Begin()
	if err != nil {
		//log.Fatal("Beginning transaction failed: " + err.Error())
		panic(err)
		//errors[0] = err
	}

	defer transaction.Rollback()
	statement, err := transaction.Prepare(sql)
	if err != nil {
		//log.Fatal("Preparing statement failed: " + err.Error())
		panic(err)
		//errors[1] = err
	}

	rows, err := statement.Query(args...)
	if err != nil {
		//log.Fatal("Running query: " + sql + " failed with: " + err.Error())
		panic(err)
		//errors[2] = err
	}

	statement.Close()
	errors[3] = transaction.Commit()
	if errors[3] != nil {
		//log.Print(errors[3])
		panic(err)
	}

	return rows, errors
}

func TransactionalExec(string string, args ...interface{}) (*sql.Result, []error) {
	return nil, nil
}

// Closes both the transaction and the statement
func CloseTransactionalStatement(statement sql.Stmt, transaction sql.Tx) error {
	return nil
}

//</editor-fold>
