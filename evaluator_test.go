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
		{
			expr: "var1 == 5.5 + 4.5 * 2",
			variables: []evaluator.Variables{
				{"var1": 14.5},
				{"var1": 20.0},
			},
			expected: []interface {
			}{
				true,
				false,
			},
		},
		{
			expr: "(5.5 + 4.5) * var1",
			variables: []evaluator.Variables{
				{"var1": 3},
				{"var1": 2},
			},
			expected: []interface {
			}{
				30,
				20,
			},
		},
		{
			expr: "as_string(coalesce(rate(var1, var2),nil,``,' '))",
			variables: []evaluator.Variables{
				{"var1": 3, "var2": 0},
				{"var1": 2, "var2": 1},
			},
			expected: []interface {
			}{
				"",
				"2",
			},
		},
		{
			expr: "coalesce(as_numeric(var1),10.0)",
			variables: []evaluator.Variables{
				{"var1": "hoge"},
				{"var1": 2.0},
				{"var1": "5.0"},
			},
			expected: []interface{}{
				10.0,
				2.0,
				5.0,
			},
		},
		{
			expr: "if(regexp_match(as_string(var1), `^hoge`), 1.8, 3.14)",
			variables: []evaluator.Variables{
				{"var1": "hoge"},
				{"var1": 2.0},
				{"var1": "5.0hoge"},
			},
			expected: []interface{}{
				1.8,
				3.14,
				3.14,
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

func TestEvaluatorVariableInvalid(t *testing.T) {

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
				"var1 variable not found",
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
		{
			expr: "var1 / var2 ",
			variables: []evaluator.Variables{
				{"var1": 1, "var2": 0},
			},
			expected: []string{
				"Eval(`var1 / var2`) divide by 0",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			e, err := evaluator.New(c.expr)
			require.NoError(t, err, "must parse success")
			e.Strict(true)
			t.Logf("%s", e)
			for i, v := range c.variables {
				_, err := e.Eval(v)
				require.Error(t, err, "must eval err, variables case %d", i)
				require.EqualError(t, err, c.expected[i], "must eval err msg match, variables case %d", i)
			}
		})
	}
}

func TestEvaluatorReservedError(t *testing.T) {
	cases := []struct {
		expr      string
		variables []evaluator.Variables
		expected  []func(error) bool
	}{
		{
			expr: "var1 / 0",
			variables: []evaluator.Variables{
				{"var1": 1},
				{"var2": 1},
			},
			expected: []func(error) bool{
				evaluator.IsDivideByZero,
				evaluator.IsVariableNotFound,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			e, err := evaluator.New(c.expr)
			require.NoError(t, err, "must parse success")
			e.Strict(true)
			t.Logf("%s", e)
			for i, v := range c.variables {
				_, err := e.Eval(v)
				require.Error(t, err, "must eval err, variables case %d", i)
				require.True(t, c.expected[i](err))
			}
		})
	}
}

func TestEvaluatorAsComparator(t *testing.T) {

	cases := map[string]bool{
		"(var1 / 2)":                            false,
		"(var1 < 2)":                            true,
		"3 < var1 > 4":                          true,
		"(var1 < var2) < bar2":                  true,
		"1 + 2 + 3":                             false,
		"coalesce(as_numeric(var1),10.0)":       false,
		"coalesce(as_numeric(var1),10.0) >= 10": true,
	}
	for expr, expected := range cases {
		t.Run(expr, func(t *testing.T) {
			e, err := evaluator.New(expr)
			require.NoError(t, err, "must parse success")
			c, ok := e.AsComparator()
			require.EqualValues(t, expected, ok)
			if expected {
				require.NotNil(t, c)
			}
		})
	}
}
