Things
==========

Style Constraints:

- The larger problem is decomposed into things that make sense for the problem domain.
- Each thing is a capsule of data that exposes procedures to the rest of the world.
- Data is never accessed directly, only through these procedures.
- Capsules can reappropriate procedures defined in other capsules.

Brief explanation of the Go implementation:

- The main function only checks for the required program arguments, and runs the `WordFrequencyController` which has the program logic.
- The program's logic is separated into 4 structs: `DataStorageManger`, `StopWordsManager`, `WordFrequencyManager`, and `WordFrequencyController`.
- Each of these structs handles a specific part of the logic as follows:
  - `DataStorageManager` handles the input file and splits it into words.
  - `StopWordsManager` handles the stop words file and checking whether a specific word is a stop word.
  - `WordFrequencyManager` handles and stores word frequencies, and can return a sorted slice of them on demand.
  - `WordFrequencyController` uses objects of the previous 3 structs to complete the term frequency task and print its output. 