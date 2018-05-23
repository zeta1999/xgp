package xgp

import (
	"math/rand"

	"github.com/MaxHalford/xgp/op"
	"github.com/MaxHalford/xgp/tree"
)

// A Mutator modifies a Tree in-place.
type Mutator interface {
	Apply(tree *tree.Tree, rng *rand.Rand)
}

// PointMutation picks one sub-Tree at random and replaces it's Operator.
type PointMutation struct {
	Weighting      Weighting
	MutateOperator func(op op.Operator, rng *rand.Rand) op.Operator
}

// Apply PointMutation.
func (mut PointMutation) Apply(tr *tree.Tree, rng *rand.Rand) {
	var f = func(tr *tree.Tree, depth int) (stop bool) {
		if rng.Float64() < mut.Weighting.apply(tr.Op) {
			tr.Op = mut.MutateOperator(tr.Op, rng)
		}
		return false
	}
	tr.Walk(f)
}

// HoistMutation selects a first sub-Tree from a Tree. It then selects a second
// sub-Tree from the first sub-Tree and replaces the first one with it. Hoist
// mutation is good for controlling bloat.
type HoistMutation struct {
	Picker Picker
}

// Apply HoistMutation.
func (mut HoistMutation) Apply(tr *tree.Tree, rng *rand.Rand) {
	// Hoist mutation only works if the height of Tree exceeds 1
	var height = tr.Height()
	if height < 1 {
		return
	}
	var (
		sub    = mut.Picker.Apply(tr, 1, tr.Height(), rng)
		subsub = mut.Picker.Apply(sub, 0, sub.Height()-1, rng)
	)
	*sub = *subsub
}

// SubtreeMutation selects a sub-Tree at random and replaces it with a new Tree.
// The new Tree has at most the same height as the selected sub-Tree.
type SubtreeMutation struct {
	NewTree   func(rng *rand.Rand) tree.Tree
	Crossover Crossover
}

// Apply SubtreeMutation.
func (mut SubtreeMutation) Apply(tr *tree.Tree, rng *rand.Rand) {
	var mutant = mut.NewTree(rng)
	mut.Crossover.Apply(tr, &mutant, rng)
}
