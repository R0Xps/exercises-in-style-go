Map Reduce
==========

Style Constraints:

- Input data is divided in blocks.
- A map function applies a given worker function to each block of data, potentially in parallel.
- A reduce function takes the results of the many worker functions and recombines them into a coherent output.

Brief explanation of the Go implementation:

- Dividing the data into blocks happens in the `partition` function, which returns a slice of strings each containing at most 200 lines from the input string.
- The worker function for the map stage is `splitWords`, which returns a slice of all non-stop words from the input string and a frequency of 1 for each of them (repeats allowed).
- The reduce function is `countWords`, which combines all the outputs of the map stage into a single map the contains every word and its total frequency, with no repeats this time.
- The map functions run in parallel so all workers can work at the same time since their data is not shared.
- Finally, after the reduce stage is done, the result map is sent to the `sorted` function which returns a slice of all words and frequencies from the map sorted in descending order by frequency. And the first 25 entries (or all entries if the slice is shorter than 25 elements) are printed.