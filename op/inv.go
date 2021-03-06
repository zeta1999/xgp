package op

import "fmt"

func safeDiv(a, b float64) float64 {
	if b == 0 {
		return 1
	}
	return a / b
}

// The Inv operator.
type Inv struct {
	Op Operator
}

// Eval computes the inverse of each value.
func (inv Inv) Eval(X [][]float64) []float64 {
	x := inv.Op.Eval(X)
	for i, xi := range x {
		x[i] = safeDiv(1, xi)
	}
	return x
}

// Arity of Inv is 1.
func (inv Inv) Arity() uint {
	return 1
}

// Operand returns Inv's operand or nil.
func (inv Inv) Operand(i uint) Operator {
	if i == 0 {
		return inv.Op
	}
	return nil
}

// SetOperand replaces Inv's operand if i is equal to 0.
func (inv Inv) SetOperand(i uint, op Operator) Operator {
	if i == 0 {
		inv.Op = op
	}
	return inv
}

// Simplify Inv.
func (inv Inv) Simplify() Operator {
	inv.Op = inv.Op.Simplify()
	switch operand := inv.Op.(type) {
	// 1 / (1 / x) = x
	case Inv:
		return operand.Op
	// 1 / a = b
	case Const:
		return Const{safeDiv(1, operand.Value)}
	}
	return inv
}

// Diff compute the following derivative: (1 / u)' = -u' / u².
func (inv Inv) Diff(i uint) Operator {
	return Inv{inv.Op.Diff(i)}
}

// Name of Inv is "inv".
func (inv Inv) Name() string {
	return "inv"
}

// String formatting.
func (inv Inv) String() string {
	return fmt.Sprintf("1/%s", parenthesize(inv.Op))
}
