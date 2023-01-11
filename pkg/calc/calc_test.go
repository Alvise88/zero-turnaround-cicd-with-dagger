package calc

import (
	"testing"

	"github.com/alvise88/zero-turnaround-cicd-with-dagger/internal/assert"
)

func TestSum(t *testing.T) {
	type testCase struct {
		first  int
		second int
		sum    int
	}

	cases := []testCase{
		{
			first:  0,
			second: 10,
			sum:    10,
		},
		{
			first:  10,
			second: 0,
			sum:    10,
		},
		{
			first:  2,
			second: 5,
			sum:    7,
		},
		{
			first:  1,
			second: 10,
			sum:    11,
		},
	}

	for _, tc := range cases {
		sum := Sum(tc.first, tc.second)

		assert.Equal(t, sum, tc.sum)
	}
}

func TestSub(t *testing.T) {
	type testCase struct {
		first  int
		second int
		sub    int
	}

	cases := []testCase{
		{
			first:  10,
			second: 0,
			sub:    10,
		},
		{
			first:  0,
			second: 10,
			sub:    -10,
		},
		{
			first:  10,
			second: 5,
			sub:    5,
		},
		{
			first:  5,
			second: 10,
			sub:    -5,
		},
	}

	for _, tc := range cases {
		sub := Sub(tc.first, tc.second)

		assert.Equal(t, sub, tc.sub)
	}
}

func TestMul(t *testing.T) {
	type testCase struct {
		first  int
		second int
		mul    int
	}

	cases := []testCase{
		{
			first:  5,
			second: 0,
			mul:    0,
		},
		{
			first:  0,
			second: 5,
			mul:    0,
		},
		{
			first:  5,
			second: 1,
			mul:    5,
		},
		{
			first:  1,
			second: 5,
			mul:    5,
		},
		{
			first:  5,
			second: 2,
			mul:    10,
		},
		{
			first:  5,
			second: 5,
			mul:    25,
		},
	}

	for _, tc := range cases {
		mul := Mul(tc.first, tc.second)

		assert.Equal(t, mul, tc.mul)
	}
}

func TestDiv(t *testing.T) {
	type testCase struct {
		first  int
		second int
		div    int
		ok     bool
	}

	cases := []testCase{
		{
			first:  10,
			second: 1,
			div:    10,
			ok:     true,
		},
		{
			first:  1,
			second: 10,
			div:    0,
			ok:     true,
		},
		{
			first:  10,
			second: 2,
			div:    5,
			ok:     true,
		},
		{
			first:  0,
			second: 10,
			div:    0,
			ok:     true,
		},
		{
			first:  10,
			second: 0,
			div:    0,
			ok:     false,
		},
	}

	for _, tc := range cases {
		div, err := Div(tc.first, tc.second)

		if !tc.ok && err == nil {
			t.Error("expected div error")
		}

		if tc.ok && err != nil {
			t.Error(err)
		}

		if tc.ok && err == nil {
			assert.Equal(t, div, tc.div)
		}
	}
}

func TestPow(t *testing.T) {
	type testCase struct {
		base     int
		exponent int
		power    int
	}

	cases := []testCase{
		{
			base:     2,
			exponent: 0,
			power:    1,
		},
		{
			base:     1,
			exponent: 0,
			power:    1,
		},
		{
			base:     3,
			exponent: 1,
			power:    3,
		},
		{
			base:     5,
			exponent: 1,
			power:    5,
		},
		{
			base:     2,
			exponent: 3,
			power:    8,
		},
		{
			base:     2,
			exponent: 2,
			power:    4,
		},
		{
			base:     2,
			exponent: -1,
			power:    0,
		},
	}

	for _, tc := range cases {
		power := Pow(tc.base, tc.exponent)

		assert.Equal(t, power, tc.power)
	}
}
