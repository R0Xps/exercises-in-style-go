Exercises in Programming Style in Go
====================================
Based on the book 'Exercises in Programming Style' by Cristina Videira Lopes

The goal of this project is to solve a simple computational task (term frequency analysis) in 7 different programming styles using the Go programming language.

### Task description:
- Given a text file, we want to display the 25 most frequent words and their frequencies, in descending order by frequency.
- The order of words that have the same frequency is not important.
- Words should be case-insensitive, and ignore stop words like 'the', 'for', etc.

### Styles used:
1. Monolithic
2. Pipeline
3. Things
4. Persistent Tables
5. Quarantine
6. Actors
7. Map Reduce

## How to run the program:
Run the docker container:
```shell
docker run -it r0xps/exercises-in-style-go:latest
```

This should open a bash session inside the container where you can run the commands provided by the program.

### Commands:
All commands follow this template (with a single exception that will be mentioned later):
```shell
style_name <stop_words_file> <input_file>
```

Every style has its own command being the style's name. For example:
```shell
actors /examples/stop_words.txt /examples/input/pride-and-prejudice.txt
```

The only exception to the template above is the persistent tables style, as that needs a database file so that is also a required argument:
```shell
persistent_tables <stop_words_file> <input_file> <database_file>
```
If the given database file exists, the program will retrieve the data stored in it instead of getting everything from the other files again.
Otherwise, a file will be created and used to store the list of words and stop words from the other two files, and it can be used to make future runs of the same input faster.

### Provided examples:
There are example input files available in the /examples directory inside the container.
These are:
- /examples/stop_words.txt - a list of stop_words and single letter words to be ignored.
- /examples/input/ - a directory containing 3 sample input files used for testing.
- /examples/output/ - a directory containing the outputs corresponding to each of the 3 input files in the previous directory (lines are sorted alphabetically for testing purposes).
