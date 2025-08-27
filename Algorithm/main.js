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
