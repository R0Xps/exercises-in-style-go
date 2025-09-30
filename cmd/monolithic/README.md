Monolithic
==========

Style Constraints:

- No named abstractions.
- No, or little, use of libraries.

Brief explanation of the Go implementation:

- First of all, arguments passed to the program are read and stored as paths for a stop words file, and an input file, respectively.
- Then those paths are used to read the files and store them as byte slices.
- Then we iterate over the stop words file byte by byte, keeping track of the start of the word and once we're no longer inside a word, we store the entire word in a string, and append it to the stop words slice.
- Next, we do the same for the input file, but instead of immediately appending the word to the slice, we look for it in the stop words slice. If it is found, we skip it and move to the next byte.
- Otherwise, we look for it in the wordFreq slice to see if it's already in there.
- If it is in the slice, the word's frequency is incremented, and it is moved up the list to its appropriate position, ensuring the slice is always sorted in descending order by frequency.
- Otherwise, it's appended to the end of the list with a frequency of 1.
- Finally, the words with the highest frequencies are printed (up to a maximum of 25 words if the slice has more than that)