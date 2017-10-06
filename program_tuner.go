package xgp

import (
	"math/rand"

	"github.com/MaxHalford/gago"
)

// A ProgramTuner optimizes a Program by tuning the Program's Constants.
type ProgramTuner struct {
	Program      *Program
	ConstValues  []float64
	ConstSetters []ConstantSetter
}

// String representation of a ProgramTuner.
func (progTuner ProgramTuner) String() string {
	return progTuner.Program.Root.String()
}

// newProgramTuner returns a ProgramTuner from a Program.
func newProgramTuner(prog *Program) ProgramTuner {
	var (
		consts       []float64
		constSetters []ConstantSetter
		addConst     = func(node *Node) {
			if c, ok := node.Operator.(Constant); ok {
				consts = append(consts, c.Value)
				constSetters = append(constSetters, node.newConstantSetter())
			}
		}
		progTuner = ProgramTuner{Program: prog.clone()}
	)
	// Extract all the Constants from the Program
	progTuner.Program.Root.RecApply(addConst)
	progTuner.ConstValues = consts
	progTuner.ConstSetters = constSetters
	return progTuner
}

// Clone a ProgramTuner.
func (progTuner ProgramTuner) clone() ProgramTuner {
	var clone = newProgramTuner(progTuner.Program)
	copy(clone.ConstValues, progTuner.ConstValues)
	return clone
}

func (progTuner *ProgramTuner) setProgConstants() {
	for i, constValue := range progTuner.ConstValues {
		progTuner.ConstSetters[i](constValue)
	}
}

func (progTuner *ProgramTuner) jitterConstants(rng *rand.Rand) {
	for i, constValue := range progTuner.ConstValues {
		progTuner.ConstValues[i] += constValue * rng.NormFloat64()
	}
}

// Implementation of the Genome interface from the gago package

// Evaluate method required to implement gago.Genome.
func (progTuner *ProgramTuner) Evaluate() float64 {
	progTuner.setProgConstants()
	return progTuner.Program.Evaluate()
}

// Mutate method required to implement gago.Genome.
func (progTuner *ProgramTuner) Mutate(rng *rand.Rand) {
	gago.MutNormalFloat64(progTuner.ConstValues, 0.8, rng)
}

// Crossover method required to implement gago.Genome.
func (progTuner *ProgramTuner) Crossover(progTuner2 gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var (
		o1     = progTuner.clone()
		o2     = progTuner.clone()
		c1, c2 = gago.CrossUniformFloat64(
			progTuner.ConstValues,
			progTuner2.(*ProgramTuner).ConstValues,
			rng,
		)
	)
	o1.ConstValues = c1
	o2.ConstValues = c2
	return gago.Genome(&o1), gago.Genome(&o2)
}

// Clone method required to implement gago.Genome.
func (progTuner ProgramTuner) Clone() gago.Genome {
	var clone = progTuner.clone()
	return &clone
}
