/*
Create a function that takes an integer input n, ranging from 2 to 50.

If the input is 4, the output must be printed like this:

* * * *
 * * *
  * *
   *

*/

package main

import (
	"fmt"
	"strings"
)

func printInvertedPyramid(n int) {
	if n < 2 || n > 50 {
		fmt.Println("Input must be between 2 - 50.")
		return
	}

	for i := n; i >= 1; i-- {
		space := strings.Repeat(" ", n-i)
		stars := strings.Repeat("* ", i)
		fmt.Println(space + stars)
	}
	fmt.Println()
}

func main() {
	printInvertedPyramid(4)
	printInvertedPyramid(5)
	printInvertedPyramid(10)
}
