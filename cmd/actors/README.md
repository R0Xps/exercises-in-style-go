Actors
==========

Style Constraints:

- The larger problem is decomposed into things that make sense for the problem domain.
- Each thing has a queue meant for other things to place messages in it.
- Each thing is a capsule of data that exposes only its ability to receive messages via the queue.
- Each thing has its own thread of execution independent of the others.

Brief explanation of the Go implementation:

- I used channels for sending messages here instead of queues, because channels were made for concurrency and act like queues. They can also be used to stop a goroutine by closing them.
- A `sync.WaitGroup` is used to make sure the program does not exit before all goroutines are done.
- The code is split into 5 parts, one main thread (the `main` function), and 4 goroutines each of which runs a different actor of the system.
- The 4 actors of the system are:
  - `DataStorageManager` handles everything related to the input file.
  - `StopWordsManager` handles everything about stop words, starting with reading them from a file, up to filtering words and only forwarding non-stop words.
  - `WordFrequencyManager` handles counting and sorting the words based on their frequencies.
  - `WordFrequencyController` acts as the driver code for the term frequency task