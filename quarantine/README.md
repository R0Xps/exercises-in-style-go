Quarantine
==========

Style Constraints:

- Core program functions have no side effects of any kind, including IO.
- All IO actions must be contained in computation sequences that are clearly separated from the pure functions.
- All sequences that have IO must be called from the main program.

Brief explanation of the Go implementation:

- 3 functions have IO interactions, `getInput`, `extractWords`, and `removeStopWords`.
- Each of these functions is a wrapper to an inner function that does the actual IO interactions needed.
- Every other function is a pure function, meaning that if it is given the exact same input, it should produce the same output every time.