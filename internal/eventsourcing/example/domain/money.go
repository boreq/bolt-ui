package domain

type Money struct {
	amount int
}

func NewMoney(amount int) Money {
	return Money{amount}
}

func (m Money) IsNegative() bool {
	return m.amount < 0
}

func (m Money) Add(o Money) Money {
	return Money{m.amount + o.amount}
}

func (m Money) Substract(o Money) Money {
	return Money{m.amount - o.amount}
}

func (m Money) Amount() int {
	return m.amount
}
