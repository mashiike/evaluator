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
		expected  []interface{}
	}{
		{
			expr: "var1",
			variables: []evaluator.Variables{
				{"var1": 1},
				{"var1": 3},
			},
			expected: []interface{}{
				1,
				3,
			},
		},
		{
			expr: "var1 <= var2",
			variables: []evaluator.Variables{
				{"var1": "abc", "var2": "def"},
				{"var1": 1, "var2": 2},
				{"var1": 3, "var2": 1},
			},
			expected: []interface{}{
				true,
				true,
				false,
			},
		},
		{
			expr: "var1 >= var2",
			variables: []evaluator.Variables{
				{"var1": "abc", "var2": "def"},
				{"var1": 1, "var2": 2},
				{"var1": 3, "var2": 1},
			},
			expected: []interface{}{
				false,
				false,
				true,
			},
		},
		{
			expr: "var1 == var2",
			variables: []evaluator.Variables{
				{"var1": "abc", "var2": "def"},
				{"var1": "abc", "var2": "abc"},
				{"var1": 1, "var2": 2},
				{"var1": 3, "var2": 3},
				{"var1": false, "var2": true},
				{"var1": false, "var2": false},
			},
			expected: []interface{}{
				false,
				true,
				false,
				true,
				false,
				true,
			},
		},
		{
			expr: "1.0 <= var1 <= 5",
			variables: []evaluator.Variables{
				{"var1": 2},
				{"var1": -2.0},
				{"var1": 10.0},
			},
			expected: []interface {
			}{
				true,
				false,
				false,
			},
		},
		{
			expr: "var1 / 2 <= 5.5 + 4.5",
			variables: []evaluator.Variables{
				{"var1": 2},
				{"var1": 30.0},
			},
			expected: []interface {
			}{
				true,
				false,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			e, err := evaluator.New(c.expr)
			require.NoError(t, err, "must parse success")
			t.Logf("%s", e)
			for i, v := range c.variables {
				actual, err := e.Eval(v)
				require.NoErrorf(t, err, "must eval success, variables case %d", i)
				require.EqualValuesf(t, c.expected[i], actual, "must eval result match, variables case %d", i)
			}
		})
	}
}

func TestEvaluatorVariableInvid(t *testing.T) {

	cases := []struct {
		expr      string
		variables []evaluator.Variables
		expected  []string
	}{
		{
			expr: "var1",
			variables: []evaluator.Variables{
				{"var2": 1},
			},
			expected: []string{
				"variable var1 is not givend",
			},
		},
		{
			expr: "var1 <= var2",
			variables: []evaluator.Variables{
				{"var1": true, "var2": true},
				{"var1": 1, "var2": "hoge"},
			},
			expected: []string{
				"Eval(`var1 <= var2`) v1[true]::bool and v2[true]::bool can not `<` comparatable",
				"Eval(`var1 <= var2`) v1[1]::int and v2[hoge]::string can not `<` comparatable",
			},
		},
		{
			expr: "var1 >= var2",
			variables: []evaluator.Variables{
				{"var1": true, "var2": true},
				{"var1": 1, "var2": "hoge"},
				{"var1": complex(1, 2), "var2": complex(1, 2)},
			},
			expected: []string{
				"Eval(`var1 >= var2`) v1[true]::bool and v2[true]::bool can not `>` comparatable",
				"Eval(`var1 >= var2`) v1[1]::int and v2[hoge]::string can not `>` comparatable",
				"Eval(`var1 >= var2`) v1[(1+2i)]::complex128 and v2[(1+2i)]::complex128 can not `>` comparatable",
			},
		},
		{
			expr: "var1 == var2",
			variables: []evaluator.Variables{
				{"var1": 1, "var2": "def"},
			},
			expected: []string{
				"Eval(`var1 == var2`) v1[1]::int and v2[def]::string can not `==` comparatable",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			e, err := evaluator.New(c.expr)
			require.NoError(t, err, "must parse success")
			t.Logf("%s", e)
			for i, v := range c.variables {
				_, err := e.Eval(v)
				require.Error(t, err, "must eval err, variables case %d", i)
				require.EqualError(t, err, c.expected[i], "must eval err msg match, variables case %d", i)
			}
		})
	}
}
