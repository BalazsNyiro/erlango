# Merged Erlango BNF definition

## The testing challenge
It is nearly impossible to test the full grammar, because the sub-elements are too wide, and I cannot generate all possible options to see the covered grammar.

A possible solution:
 - create small separated files and define LOCAL symbols for testing purpose
 - the small files are MERGED later, without the LOCAL sections, so the small tested units can be merged to a bigger one.

