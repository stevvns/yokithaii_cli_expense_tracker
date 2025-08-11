package operator

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
)

// Describes each expense
type expense struct {
	id     int
	date   time.Time
	desc   string
	amount int
}

// Set of getters
func (e *expense) Id() int {
	return e.id
}
func (e *expense) Amount() int {
	return e.amount
}
func (e *expense) Date() time.Time {
	return e.date
}
func (e *expense) Desc() string {
	return e.desc
}

// If provided data is valid - returns nil, else - returns error
func ValidateData(description string, amount int) error {

	if amount < 0 {
		return errors.New("negative amount promoted")
	}
	if amount > math.MaxInt {
		return errors.New("amount is too big")
	}
	runes := []rune(description)
	if len(runes) > 30 { // i do believe that u can describe your purchase in 30 symbols or less
		return errors.New("your description is too long")
	}

	return nil

}

func CreateExpense(description string, amount int) expense {
	return expense{0, time.Now(), description, amount}
}

// Shows list of available commands
func Help() {
	fmt.Println("Доступные команды: 1. add *Description* *amount* -- adds your expense")
	fmt.Println("Доступные команды: 2. list -- shows list of your expenses")
	fmt.Println("Доступные команды: 3. summary -- shows amount you spent")
	fmt.Println("Доступные команды: 4. help -- shows a list of commands")
	fmt.Println("Доступные команды: 5. delete *id* -- deletes expense with given id(rest of id's stays the same)")
	fmt.Println("Доступные команды: 6. update *id* *description* *amount* -- updates id if exists")
	fmt.Println("Доступные команды: 7. exit -- ends a program")
	fmt.Println("Доступные команды: 8. reset -- resets database")
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WORKS ONLY WITH LOCAL DB!
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Returns a connection to local DB
func NewConnectionToDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:Asdewq232113@localhost:5433/postgres")

	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Checks that expense if provided id exists
func ExpenseExists(conn *pgx.Conn, id int) (bool, error) {
	var exists bool
	err := conn.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM public.Expenses WHERE id = $1)`, id).Scan(&exists)
	return exists, err
}

// Adds a promoted expense (need to validate)
func AddExpensToDatabase(conn *pgx.Conn, e expense) error {
	query := `INSERT INTO public.Expenses (Date, Description, Amount) VALUES (@Date, @Description, @Amount)`

	args := pgx.NamedArgs{
		"Date":        e.date,
		"Description": e.desc,
		"Amount":      e.amount,
	}

	_, err := conn.Exec(context.Background(), query, args)
	if err != nil {
		return fmt.Errorf("error inserting expense: %v", err)
	}
	return nil
}

// Returns an array of all expenses available
func GetAllExpensesFromDb(conn *pgx.Conn) (*[]expense, error) {

	query := `SELECT * FROM public.Expenses`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error querying the table")
		return nil, err
	}

	defer rows.Close()

	var expenses []expense

	for rows.Next() {
		var exp expense
		err := rows.Scan(&exp.id, &exp.date, &exp.desc, &exp.amount)
		if err != nil {
			fmt.Println("Error fetching expense details")
			return &expenses, err
		}
		expenses = append(expenses, exp)
	}
	return &expenses, nil
}

// Replaces expense with promoted id with new given expense
func UpdateExpense(conn *pgx.Conn, expId int, exp expense) error {
	query := `
		UPDATE public.Expenses
		SET Date = @Date, Description = @Description, Amount = @Amount 
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":          expId,
		"Date":        exp.date,
		"Description": exp.desc,
		"Amount":      exp.amount,
	}

	_, err := conn.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}
	return nil
}

// Deletes an expense with promoted id
func DeleteExpense(conn *pgx.Conn, id int) error {

	query := `
	DELETE FROM public.Expenses WHERE id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := conn.Exec(context.Background(), query, args)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// to help count sum of array
func sumOfArray(i []int) int {
	summ := 0
	for _, v := range i {
		summ += v
	}
	return summ
}

// Returns an amount of all expenses
func Summary(conn *pgx.Conn) (int, error) {

	query := `
	SELECT e.Amount FROM public.Expenses as e
	`

	rows, err := conn.Query(context.Background(), query)

	if err != nil {
		fmt.Println("Error getting summary")
		return 0, err
	}
	defer rows.Close()

	var amounts []int

	for rows.Next() {
		var temp int
		err := rows.Scan(&temp)
		if err != nil {
			fmt.Println("Error getting one amount")
			return sumOfArray(amounts), err
		}
		amounts = append(amounts, temp)
	}
	return sumOfArray(amounts), nil

}

// Wipes out db (and resets id's too!)
func ResetDatabase(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(),
		`TRUNCATE TABLE public.Expenses RESTART IDENTITY`)
	return err
}
