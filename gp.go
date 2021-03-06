package xgp

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/MaxHalford/eaopt"
	"github.com/gosuri/uiprogress"

	"github.com/MaxHalford/xgp/metrics"
	"github.com/MaxHalford/xgp/op"
)

// An GP links all the different components together and can be used to train
// Programs on a dataset. You shouldn't instantiate this struct directly;
// instead you should use the GPConfig struct and call it's NewGP method.
type GP struct {
	GPConfig

	EvalMetric       metrics.Metric
	LossMetric       metrics.Metric
	Functions        []op.Operator
	Initializer      Initializer
	GA               *eaopt.GA
	PointMutation    PointMutation
	SubtreeMutation  SubtreeMutation
	HoistMutation    HoistMutation
	SubtreeCrossover SubtreeCrossover

	fm       map[uint][]op.Operator
	X        [][]float64
	Y        []float64
	W        []float64
	XVal     [][]float64
	YVal     []float64
	WVal     []float64
	nClasses int
}

// String representation of an GP.
func (gp GP) String() string {
	return gp.GPConfig.String()
}

// BestProgram returns the GP's best obtained Program.
func (gp GP) BestProgram() (Program, error) {
	if len(gp.GA.HallOfFame) == 0 {
		return Program{}, errors.New("The GP has not been trained yet")
	}
	return *gp.GA.HallOfFame[0].Genome.(*Program), nil
}

func (gp GP) progress(start time.Time) string {
	// Add time spent
	var message = fmtDuration(time.Since(start))
	// Add training error
	best, err := gp.BestProgram()
	if err != nil {
		return ""
	}
	yTrainPred, err := best.Predict(gp.X, gp.EvalMetric.NeedsProbabilities())
	if err != nil {
		return ""
	}
	trainScore, err := gp.EvalMetric.Apply(gp.Y, yTrainPred, nil)
	if err != nil {
		return ""
	}
	message += fmt.Sprintf(", train %s: %.5f", gp.EvalMetric.String(), trainScore)
	// Add validation error
	if gp.XVal != nil && gp.YVal != nil {
		yEvalPred, err := best.Predict(gp.XVal, gp.EvalMetric.NeedsProbabilities())
		if err != nil {
			return ""
		}
		evalScore, err := gp.EvalMetric.Apply(gp.YVal, yEvalPred, gp.WVal)
		if err != nil {
			return ""
		}
		message += fmt.Sprintf(", val %s: %.5f", gp.EvalMetric.String(), evalScore)
	}
	return message
}

// polishBest takes the best Program and polishes it.
func (gp *GP) polishBest() error {
	best, err := gp.BestProgram()
	if err != nil {
		return err
	}
	polished, err := polishProgram(best, gp.RNG)
	if err != nil {
		return err
	}
	fitness, err := polished.Evaluate()
	if err != nil {
		return err
	}
	if fitness < gp.GA.HallOfFame[0].Fitness {
		gp.GA.HallOfFame[0].Genome = &polished
	}
	return nil
}

// Fit an GP to a dataset.
func (gp *GP) Fit(
	// Required arguments
	X [][]float64,
	Y []float64,
	// Optional arguments (can safely be nil)
	W []float64,
	XVal [][]float64,
	YVal []float64,
	WVal []float64,
	verbose bool,
) error {

	// Set the training set
	gp.X = X
	gp.Y = Y
	gp.W = W

	// Set the validation set
	gp.XVal = XVal
	gp.YVal = YVal
	gp.WVal = WVal

	// Count the number of classes if the task is classification
	if gp.LossMetric.Classification() {
		// Check that the task to perform is not multi-class classification
		gp.nClasses = countDistinct(Y)
		if gp.nClasses > 2 {
			return errors.New("Multi-class classification is not supported")
		}
	}

	// Evolve the GA
	var (
		bar      *uiprogress.Bar
		progress *uiprogress.Progress
	)
	if verbose {
		// Initialize a progress bar
		var start = time.Now()
		progress = uiprogress.New()
		progress.Start()
		var steps = int(gp.NGenerations)
		if gp.PolishBest {
			steps++
		}
		bar = progress.AddBar(steps)
		bar.PrependCompleted()
		bar.AppendFunc(func(b *uiprogress.Bar) string {
			return gp.progress(start)
		})
		// Make sure the progress bar will stop
		defer func() { progress.Stop() }()
		// Use a callback to increment the progress bar at each generation
		gp.GA.Callback = func(ga *eaopt.GA) { bar.Incr() }
	}

	// Run the GA
	err := gp.GA.Minimize(func(rng *rand.Rand) eaopt.Genome {
		var prog = gp.newProgram(rng)
		return &prog
	})
	if err != nil {
		return err
	}

	// Polish the best Program
	if gp.PolishBest {
		err := gp.polishBest()
		if err != nil {
			return err
		}
		if verbose {
			bar.Incr()
		}
	}

	return nil
}

// Predict makes predictions with the best obtained Program as so far.
func (gp GP) Predict(X [][]float64, proba bool) ([]float64, error) {
	var best, err = gp.BestProgram()
	if err != nil {
		return nil, err
	}
	return best.Predict(X, proba)
}

// PredictPartial is a convenience function on top of Predict to make
// predictions on a single instance.
func (gp GP) PredictPartial(x []float64, proba bool) (float64, error) {
	var best, err = gp.BestProgram()
	if err != nil {
		return 0, err
	}
	return best.PredictPartial(x, proba)
}

func (gp GP) newConst(rng *rand.Rand) op.Const {
	return op.Const{
		Value: gp.GPConfig.ConstMin + rng.Float64()*(gp.GPConfig.ConstMax-gp.ConstMin),
	}
}

func (gp GP) newVar(rng *rand.Rand) op.Var {
	return op.Var{Index: uint(rng.Intn(len(gp.X)))}
}

func (gp GP) newFunction(rng *rand.Rand) op.Operator {
	return gp.Functions[rng.Intn(len(gp.Functions))]
}

func (gp GP) newFunctionOfArity(arity uint, rng *rand.Rand) op.Operator {
	n := len(gp.fm[arity])
	if n == 0 {
		return nil
	}
	return gp.fm[arity][rng.Intn(n)]
}

func (gp GP) newOperator(rng *rand.Rand) op.Operator {
	return gp.Initializer.Apply(
		gp.MinHeight,
		gp.MaxHeight,
		func(leaf bool, rng *rand.Rand) op.Operator {
			if leaf {
				if rng.Float64() < gp.PConst {
					return gp.newConst(rng)
				}
				return gp.newVar(rng)
			}
			return gp.newFunction(rng)
		},
		rng,
	)
}

func (gp GP) newProgram(rng *rand.Rand) Program {
	return Program{
		Op: gp.newOperator(rng),
		GP: &gp,
	}
}

func (gp GP) mutateOperator(operator op.Operator, rng *rand.Rand) op.Operator {
	switch operator.(type) {
	case op.Const:
		return op.Const{Value: operator.(op.Const).Value * rng.NormFloat64()}
	case op.Var:
		return gp.newVar(rng)
	default:
		newOp := gp.newFunctionOfArity(operator.Arity(), rng)
		// newFunctionOfArity might return nil if there are no available
		// operators of the given arity
		if newOp == nil {
			return operator
		}
		// Don't forget to set the new Operator's operands
		for i := uint(0); i < operator.Arity(); i++ {
			newOp = newOp.SetOperand(i, operator.Operand(i))
		}
		return newOp
	}
}
