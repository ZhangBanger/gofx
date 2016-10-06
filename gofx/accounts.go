package gofx

type Account struct {
	User    string  `db:"id"`
	Balance float64 `db:"balance"`
}
