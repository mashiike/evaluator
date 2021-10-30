package evaluator_test

import (
	"testing"

	"github.com/mashiike/evaluator"
	"github.com/stretchr/testify/require"
)

func TestEvaluatorSuccess(t *testing.T) {

	cases := []struct {
		expr      string
		variables []evaluator.Variables
		expected  []bool
	}{
		{
			expr: "var1 <= var2",
			variables: []evaluator.Variables{
				{"var1": 1, "var2": 2},
				{"var1": 3, "var2": 1},
			},
			expected: []bool{
				true,
				false,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			e, err := evaluator.New(c.expr)
			require.NoError(t, err, "must parse success")
			for i, v := range c.variables {
				actual, err := e.Eval(v)
				require.NoErrorf(t, err, "must eval sucess, variables case %d", i)
				require.EqualValuesf(t, c.expected[i], actual, "must eval result match, variables case %d", i)
			}
		})
	}

}
