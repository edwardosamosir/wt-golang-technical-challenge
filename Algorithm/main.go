/*
	Using your preferred programming language, develop the most efficient function that returns a
	list of all possible combinations of non-repeating digit (1-9) given variable l and t. l is the length
	of a combination and t is the total of all numbers in the combination.

	Rules:
	- Digit ranges from 1 to 9
	- A digit can only be used once
	- Combination must not appear twice. Example: [1,2,3] is the same with [3,2,1] so only one of
	them should be in the list.

	Example 1:
	Input: l= 3, t = 6
	Output: [[1,2,3]]
	Explanation:
	1 + 2 + 3 = 6

	Example 2:
	Input: l = 3, t = 8
	Output: [[1,2,5], [1,3,4]]
	Explanation:
	1 + 2 + 5 = 8 ,1 + 3 + 4 = 8

	Example 3:
	Input: l = 4, t = 5
	Output: []
	Explanation: no combination
*/

package main

import "fmt"

func combinationSum(l int, t int) [][]int {
	result := [][]int{}
	path := []int{}

	var backtrack func(start, remain, remainSum int)
	backtrack = func(start, remain, remainSum int) {
		if remain == 0 && remainSum == 0 {
			comb := append([]int{}, path...)
			result = append(result, comb)
			return
		}

		if remain < 0 || remainSum < 0 {
			return
		}

		if start > 10-remain {
			return
		}

		minSum := remain*start + (remain*(remain-1))/2

		maxSum := remain * (19 - remain) / 2
		if remainSum < minSum || remainSum > maxSum {
			return
		}

		for num := start; num <= 9; num++ {
			if num > remainSum {
				break
			}

			path = append(path, num)
			backtrack(num+1, remain-1, remainSum-num)
			path = path[:len(path)-1]
		}
	}

	backtrack(1, l, t)
	return result
}

func main() {
	// Example 1
	fmt.Println("Example 1:", combinationSum(3, 6))
	// Example 2
	fmt.Println("Example 2:", combinationSum(3, 8))
	// Example 3
	fmt.Println("Example 3:", combinationSum(4, 5))
	// Example 4
	fmt.Println("Example 4:", combinationSum(5, 30))
	// Example 5
	fmt.Println("Example 5:", combinationSum(7, 40))
	// Example 6
	fmt.Println("Example 6:", combinationSum(9, 45))
}
