Pipeline
==========

Style Constraints:

- Larger problem is decomposed using functional abstraction. Functions take input and produce output.
- No shared state between functions.
- The larger problem is solved by composing functions one after the other, in pipeline, as a faithful reproduction of mathematical function composition f ◦ g (f after g).

Brief explanation of the Go implementation:

- In the pipeline style, all operations are split into functions executed in sequence.
- And in cases where more than 1 function parameter is necessary, currying can be used to convert it into a sequence of functions that take a single argument each.
- The order of operations (and function calls) is as follows:
  1. Read the input file from the path given as an argument to the program.
  2. Filter the file's contents and normalize them to be all lowercase letters and spaces only.
  3. Split the filtered string into a slice of all the words in it.
  4. Remove all the stop words (which are read from a file in the other path given to the program as an argument) from the words slice.
  5. Make a map containing all the words and their frequencies.
  6. Copy the contents of the map into a slice that's sorted by frequency in descending order.
  7. Print the first 25 elements from the final words slice (or all of the elements if the slice contains less than 25 elements).