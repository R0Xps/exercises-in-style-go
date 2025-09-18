Persistent Tables
==========

Style Constraints:

- The data exists beyond the execution of programs that use it, and is
  meant to be used by many different programs.
- The data is stored in a way that makes it easier/faster to explore. For
  example:
  - The input data of the problem is modeled as one or more series of
  domains, or types, of data.
  - The concrete data is modeled as having components of several
  domains, establishing relationships between the application’s data
  and the domains identified.
- The problem is solved by issuing queries over the data.

Brief explanation of the Go implementation:

- The code of this style requires an additional command-line argument that is the database file path.
- If the given file exists, an sqlite database is read from it and used to get the word count.
- And if it doesn't exist, the tables are created (the file is automatically created in the process), and the stop words and input files are inserted into the appropriate tables.
- Then a database query gets a list of the words with the most frequencies (max 25 words) and the results are printed the same as the other styles.