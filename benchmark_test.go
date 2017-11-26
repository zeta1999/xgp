package koza

import (
	"testing"

	"github.com/MaxHalford/koza/dataset"
)

func BenchmarkFit(b *testing.B) {
	var estimator, err = NewEstimatorWithDefaults()
	if err != nil {
		panic(err)
	}
	// Load the training set in memory
	boston, err := dataset.ReadCSV("examples/boston/train.csv", "y", false)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		estimator.Fit(boston.X, boston.Y, boston.XNames, false)
	}
}