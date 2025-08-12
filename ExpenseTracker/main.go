package main

import (
	"ExpenseTracker/operator"
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	conn, err := operator.NewConnectionToDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	scanner := bufio.NewScanner(os.Stdin)

	for {
		operator.Help()
		fmt.Print("Введите команду: ")

		if ok := scanner.Scan(); !ok {
			fmt.Println("Ошибка ввода!")
			operator.Help()
		}
		text := scanner.Text()
		fields := strings.Fields(text)

		if len(fields) == 0 {
			fmt.Println("Вы ничего не ввели")
			operator.Help()
		}

		cmd := fields[0]

		if cmd == "exit" {
			return

		} else if cmd == "add" {
			if len(fields) < 3 {
				fmt.Println("Not enough arguments provided for command \"add\"")
				operator.Help()
			}
			temp, err := strconv.Atoi(fields[2])
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = operator.ValidateData(fields[1], temp)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = operator.AddExpensToDatabase(conn, operator.CreateExpense(fields[1], temp))
			if err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Println("Expense added succesfully")
			}

		} else if cmd == "list" {
			expenses, err := operator.GetAllExpensesFromDb(conn)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%-4s %-12s %-20s %s\n", "#", "Date", "Description", "Amount")
			fmt.Println("--------------------------------------------------")
			for _, v := range *expenses {
				fmt.Printf("%-4v %-12v %-20v $%d\n",
					v.Id(),
					v.Date().Format("2006-01-02"),
					v.Desc(),
					v.Amount())
			}

		} else if cmd == "summary" {
			summary, err := operator.Summary(conn)
			if err != nil {
				fmt.Println(err)
				continue
			} else {
				fmt.Println(summary)
			}

		} else if cmd == "help" {
			operator.Help()

		} else if cmd == "delete" {
			if len(fields) < 2 {
				fmt.Println("You provided no id to delete")
				operator.Help()
				continue
			}
			identif, err := strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println(err)
			} else {
				exists, err := operator.ExpenseExists(conn, identif)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if !exists {
					fmt.Println("No expense with given id")
					continue
				}
				err = operator.DeleteExpense(conn, identif)
				if err != nil {
					fmt.Println(err)
					continue
				} else {
					fmt.Printf("Expense deleted succesfully (ID: %v)\n", identif)
				}

			}

		} else if cmd == "update" {
			if len(fields) < 4 {
				fmt.Println("Not enough arguments provided for command \"update\"")
				operator.Help()
				continue
			}

			id, err := strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println(err)
				continue
			}

			exists, err := operator.ExpenseExists(conn, id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if !exists {
				fmt.Println("No expense with provided id")
				continue
			}

			amount, err := strconv.Atoi(fields[3])
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = operator.ValidateData(fields[2], amount)
			if err != nil {
				fmt.Println(err)
				continue
			}

			updatedExpense := operator.CreateExpense(fields[2], amount)

			err = operator.UpdateExpense(conn, id, updatedExpense)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Succesfully updated ID: %v\n", id)
			}
		} else if cmd == "reset" {
			err = operator.ResetDatabase(conn)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			fmt.Println("No such coommand!")
			operator.Help()
		}

	}

}
