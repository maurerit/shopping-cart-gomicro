package repository

import (
	"database/sql"
	"log"
)

const (
	SHOPPING            = 1 << iota
	ORDERED             = 1 << iota
	IN_PROGRESS         = 1 << iota
	GATHERING_MATERIALS = 1 << iota
	SHIPPING            = 1 << iota
	COMPLETE            = 1 << iota
)

// Global database handle
var DB *sql.DB

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
