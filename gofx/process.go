package gofx

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math"
)

func Process() {
	defer dbMap.Db.Close()

	for order := range orderChan {
		// Get possible transactions
		var matchingOrders []Order

		if order.Buy {
			dbMap.Select(
				&matchingOrders,
				"select * from orders where Security=? and Buy=0 and Price < ? order by ts",
				order.Security,
				order.Price,
			)
		} else {
			dbMap.Select(
				&matchingOrders,
				"select * from orders where Security=? and Buy=1 and Price > ? order by ts",
				order.Security,
				order.Price,
			)
		}

		for _, matchingOrder := range matchingOrders {
			// Calculate diff for this potential transaction
			amountToDeduct := minInt(order.Quantity, matchingOrder.Quantity)
			if amountToDeduct == 0 {
				break
			}
			priceDiff := math.Abs(order.Price - matchingOrder.Price)

			// Apply deductions to matching
			matchingOrder.Quantity -= amountToDeduct
			order.Quantity -= amountToDeduct
			if matchingOrder.Quantity == 0 {
				dbMap.Delete(&matchingOrder)
			} else {
				dbMap.Update(&matchingOrder)
			}

			// Update user accounts
			updateUserAccount(
				dbMap,
				order.User,
				!order.Buy,
				float64(amountToDeduct)*priceDiff,
			)
			updateUserAccount(
				dbMap,
				matchingOrder.User,
				order.Buy,
				float64(amountToDeduct)*priceDiff,
			)
		}

		// If order still has fill, write to db
		if order.Quantity > 0 {
			dbMap.Insert(&order)
		}
	}
}

func minInt(left uint32, right uint32) uint32 {
	if left < right {
		return left
	} else {
		return right
	}
}

func updateUserAccount(dbMap *gorp.DbMap, user string, add bool, amount float64) {
	var account *Account

	accInter, _ := dbMap.Get(Account{}, user)
	if accInter == nil {
		account = &Account{user, 0.0}
		fmt.Println("creating user:", *account)
		err := dbMap.Insert(&account)
		checkErr(err, "Insert failed for account"+account.User)
	} else {
		account = accInter.(*Account)
	}

	if add {
		account.Balance += amount
	} else {
		account.Balance -= amount
	}

	dbMap.Update(&account)
}
