package calc

import "errors"

func Sum(first, second int) int {
	a := first
	b := second

	for b != 0 {
		carry := a & b // Carry value is calculated
		a ^= b         // Sum value is calculated and stored in a
		b = carry << 1 // The carry value is shifted towards left by a bit
	}

	return a // returns the final sum
}

func Sub(first, second int) int {
	return Sum(first, -second)
}

func Mul(first, second int) int {
	var mul int
	for i := 0; i < second; i++ {
		mul = Sum(mul, first)
	}
	return mul
}

func Div(first, second int) (int, error) {
	if second == 0 {
		return 0, errors.New("division by zero")
	}

	if second > first {
		return 0, nil
	}

	var div int

	rest := first
	for rest > 0 {
		div++

		rest -= second
	}

	return div, nil
}

func Pow(base, exponent int) int {
	if exponent < 0 {
		return 0
	}

	if exponent == 0 {
		return 1
	}

	if exponent == 1 {
		return base
	}

	return base * Pow(base, exponent-1)
}
