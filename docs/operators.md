# Available operators

The following table lists all the available operators. Regardless of from where it is being used from, functions should be passed to koza by concatenating the short names of the functions with a comma. For example to use the natural logarithm and the product one should pass `log,mul` as an argument.

Code-wise the operators are all located in the `op` sub-package, of which the goal is to provide fast implementations for each operator. For the while the only accelerations that exist are the ones for the sum and the division which use assembly implementations made available by [gonum/floats](https://godoc.org/gonum.org/v1/gonum/floats).

| Name | Arity | Short name | Golang struct | Assembly code |
|------|-------|------------|---------------|---------------|
| Cosine | 1 | cos | Cos | ✗ |
| Sine | 1 | sin | Sin | ✗ |
| Natural logarithm | 1 | log | Log | ✗ |
| Exponential | 1 | exp | Exp | ✗ |
| Maximum | 2 | max | Max | ✗ |
| Minimum | 2 | min | Min | ✗ |
| Sum | 2 | sum | Sum | ✔ |
| Difference | 2 | sub | Difference | ✗ |
| Division | 2 | div | Division | ✔ |
| Product | 2 | mul | Product | ✗ |
| Power | 2 | pow | Power | ✗ |