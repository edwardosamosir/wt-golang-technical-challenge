/*
Create a function that takes an integer input n, ranging from 2 to 50.

If the input is 4, the output must be printed like this:

* * * *
 * * *
  * *
   *

*/
function printInvertedPyramid(n) {
  if (n < 2 || n > 50) {
    console.log("Input must be between 2 - 50.");
    return;
  }

  for (let i = n; i >= 1; i--) {
    const spaces = " ".repeat(n - i);
    const stars = "* ".repeat(i);
    console.log(spaces + stars);
  }
  console.log();
}


printInvertedPyramid(4);
printInvertedPyramid(5);
printInvertedPyramid(10);
