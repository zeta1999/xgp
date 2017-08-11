package xgp

import (
	"math"
)

// FUNCTIONS maps Function string representation to Functions for serialization
// purposes.
var FUNCTIONS = map[string]Operator{
	Cos{}.String():        Cos{},
	Sin{}.String():        Sin{},
	Log{}.String():        Log{},
	Exp{}.String():        Exp{},
	Max{}.String():        Max{},
	Min{}.String():        Min{},
	Sum{}.String():        Sum{},
	Difference{}.String(): Difference{},
	Division{}.String():   Division{},
	Product{}.String():    Product{},
	Power{}.String():      Power{},
}

// 1D functions

// Cos computes the cosine of an operand.
type Cos struct{}

// Apply Cos.
func (op Cos) Apply(X []float64) float64 {
	return math.Cos(X[0])
}

// Arity of Cos.
func (op Cos) Arity() int {
	return 1
}

// String representation of Cos.
func (op Cos) String() string {
	return "cos"
}

// Sin computes the sine of an operand.
type Sin struct{}

// Apply Sin.
func (op Sin) Apply(X []float64) float64 {
	return math.Sin(X[0])
}

// Arity of Sin.
func (op Sin) Arity() int {
	return 1
}

// String representation of Sin.
func (op Sin) String() string {
	return "sin"
}

// Log computes the natural logarithm of an operand.
type Log struct{}

// Apply Log.
func (op Log) Apply(X []float64) float64 {
	return math.Log(X[0])
}

// Arity of Log.
func (op Log) Arity() int {
	return 1
}

// String representation of Log.
func (op Log) String() string {
	return "log"
}

// Exp computes the exponential of an operand.
type Exp struct{}

// Apply Exp.
func (op Exp) Apply(X []float64) float64 {
	return math.Exp(X[0])
}

// Arity of Exp.
func (op Exp) Arity() int {
	return 1
}

// String representation of Exp.
func (op Exp) String() string {
	return "exp"
}

// 2D operators

// Max returns the maximum of two operands.
type Max struct{}

// Apply Max.
func (op Max) Apply(X []float64) float64 {
	if X[0] > X[1] {
		return X[0]
	}
	return X[1]
}

// Arity of Max.
func (op Max) Arity() int {
	return 2
}

// String representation of Max.
func (op Max) String() string {
	return "max"
}

// Min returns the minimum of two operands.
type Min struct{}

// Apply Min.
func (op Min) Apply(X []float64) float64 {
	if X[0] < X[1] {
		return X[0]
	}
	return X[1]
}

// Arity of Min.
func (op Min) Arity() int {
	return 2
}

// String representation of Min.
func (op Min) String() string {
	return "min"
}

// Sum returns the sum of two operands.
type Sum struct{}

// Apply Sum.
func (op Sum) Apply(X []float64) float64 {
	return X[0] + X[1]
}

// Arity of Sum.
func (op Sum) Arity() int {
	return 2
}

// String representation of String.
func (op Sum) String() string {
	return "+"
}

// Difference returns the difference between two operands.
type Difference struct{}

// Apply Difference.
func (op Difference) Apply(X []float64) float64 {
	return X[0] - X[1]
}

// Arity of Difference.
func (op Difference) Arity() int {
	return 2
}

// String representation of Difference.
func (op Difference) String() string {
	return "-"
}

// Division returns the division of two operands. The left operand is the
// numerator and the right operand is the denominator. The division is protected
// so that if the denominator's value is in range [-0.001, 0.001] the operator
// returns 1.
type Division struct{}

// Apply Division.
func (op Division) Apply(X []float64) float64 {
	if math.Abs(X[1]) < 0.001 {
		return 1
	}
	return X[0] / X[1]
}

// Arity of Division.
func (op Division) Arity() int {
	return 2
}

// String representation of Division.
func (op Division) String() string {
	return "/"
}

// Product returns the product two operands.
type Product struct{}

// Apply Product.
func (op Product) Apply(X []float64) float64 {
	return X[0] * X[1]
}

// Arity of Product.
func (op Product) Arity() int {
	return 2
}

// String representation of Product.
func (op Product) String() string {
	return "*"
}

// Power computes the exponent of a first value by a second one.
type Power struct{}

// Apply Power.
func (op Power) Apply(X []float64) float64 {
	return math.Pow(X[0], X[1])
}

// Arity of Power.
func (op Power) Arity() int {
	return 2
}

// String representation of Power.
func (op Power) String() string {
	return "^"
}
