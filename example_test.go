package evaluator_test

import (
	"fmt"
	"log"

	"github.com/mashiike/evaluator"
)

func ExampleEvaluator_Eval() {

	e, _ := evaluator.New("(var1 + 0.5) * var2")
	ans, _ := e.Eval(evaluator.Variables{"var1": 0.5, "var2": 3})
	fmt.Println(ans)

	// Output:
	// 3
}

func ExampleComparator_Compare() {

	e, err := evaluator.New("(var1 + 0.5) <= var2")
	if err != nil {
		log.Fatal(err)
	}
	c, ok := e.AsComparator()
	if !ok {
		log.Fatal("not comparative expr")
	}
	ans, err := c.Compare(evaluator.Variables{"var1": 0.5, "var2": 3})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ans)

	// Output:
	// true
}
