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

function combinationSum(l, t) {
  const result = [];
  const path = [];

  function backtrack(start, remain, remainSum) {
    if (remain === 0 && remainSum === 0) {
      result.push([...path]); 
      return;
    }

    if (remain < 0 || remainSum < 0) {
      return;
    }

    if (start > 10 - remain) {
      return;
    }

    const minSum = remain * start + (remain * (remain - 1)) / 2;

    const maxSum = (remain * (19 - remain)) / 2;

    if (remainSum < minSum || remainSum > maxSum) {
      return;
    }

    for (let num = start; num <= 9; num++) {
      if (num > remainSum) break;

      path.push(num);
      backtrack(num + 1, remain - 1, remainSum - num);
      path.pop();
    }
  }

  backtrack(1, l, t);
  return result;
}

// Test cases
console.log("Example 1:", combinationSum(3, 6));
console.log("Example 2:", combinationSum(3, 8));
console.log("Example 3:", combinationSum(4, 5));
console.log("Example 4:", combinationSum(5, 30));
console.log("Example 5:", combinationSum(7, 40));
console.log("Example 6:", combinationSum(9, 45));
