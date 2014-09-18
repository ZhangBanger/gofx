package gofx

import (
	_ "github.com/go-sql-driver/mysql"
)

func Process() {
	defer dbMap.Db.Close()

	for order := range orderChan {
		// Get matching orders
		var matchingOrders []Order
		if order.Buy {
			dbMap.Select(
				&matchingOrders,
				"select * from orders where security=? and buy=0 and price < ? order by timestamp",
				order.Security,
				order.Price,
			)
		} else {
			dbMap.Select(
				&matchingOrders,
				"select * from orders where security=? and buy=1 and price > ? order by timestamp",
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

			// Apply deductions to matching
			matchingOrder.Quantity -= amountToDeduct
			order.Quantity -= amountToDeduct

			if matchingOrder.Quantity == 0 {
				dbMap.Delete(&matchingOrder)
			} else {
				dbMap.Update(&matchingOrder)
			}
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
