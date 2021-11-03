package evaluator_test

import (
	"math/rand"
	"testing"

	"github.com/mashiike/evaluator"
	"github.com/stretchr/testify/require"
)

func BenchmarkEvalutorEval(b *testing.B) {

	cases := []struct {
		casename      string
		variablesFunc func() evaluator.Variables
		expr          string
	}{
		{
			casename: "ref_only",
			variablesFunc: func() evaluator.Variables {
				return evaluator.Variables{
					"var1": rand.NormFloat64(),
				}
			},
			expr: "var1",
		},
		{
			casename: "simple",
			variablesFunc: func() evaluator.Variables {
				return evaluator.Variables{
					"var1": rand.NormFloat64(),
				}
			},
			expr: "var1 <= 30.0",
		},
		{
			casename: "add_compare",
			variablesFunc: func() evaluator.Variables {
				return evaluator.Variables{
					"var1": rand.NormFloat64(),
					"var2": rand.NormFloat64(),
				}
			},
			expr: "var1 + var2 <= 30.0",
		},
		{
			casename: "rate",
			variablesFunc: func() evaluator.Variables {
				return evaluator.Variables{
					"var1": rand.NormFloat64(),
					"var2": rand.NormFloat64(),
				}
			},
			expr: "rate(var1,var1 + var2) <= 0.95",
		},
		{
			casename: "full",
			variablesFunc: func() evaluator.Variables {
				return evaluator.Variables{
					"var1": randStringVariable(),
					"var2": rand.NormFloat64(),
				}
			},
			expr: "(coalesce(as_numeric(var1),10.0) + 300.0) * var2 <= 30.0",
		},
	}
	for _, c := range cases {
		e, err := evaluator.New(c.expr)
		require.NoError(b, err)
		b.Run(c.casename, func(b *testing.B) {
			varsSlice := make([]evaluator.Variables, 0, b.N)
			for i := 0; i < b.N; i++ {
				varsSlice = append(varsSlice, c.variablesFunc())
			}
			b.ResetTimer()
			for _, vars := range varsSlice {
				e.Eval(vars)
			}
		})
	}
}

var numOnlyLetters = []rune("0123456789")
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randStringVariable() interface{} {
	r := rand.Intn(100)
	if r < 30 {
		return rand.NormFloat64()
	} else if r < 60 {
		return randString(10, letters)
	} else {
		return randString(10, numOnlyLetters)
	}
}

func randString(n int, l []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = l[rand.Intn(len(l))]
	}
	return string(b)
}
