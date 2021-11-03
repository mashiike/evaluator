package evaluator_test

import (
	"fmt"

	"github.com/mashiike/evaluator"
)

func ExampleEvaluator_Eval() {

	e, _ := evaluator.New("(var1 + 0.5) * var2")
	ans, _ := e.Eval(evaluator.Variables{"var1": 0.5, "var2": 3})
	fmt.Println(ans)

	// Output:
	// 3
}
