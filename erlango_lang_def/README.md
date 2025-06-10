# Merged Erlango BNF definition
The grammar is stored in ordered .bnf files.
Those are converted/merged for LARK. 

## The testing challenge
It is nearly impossible to test the full grammar, because the sub-elements are too wide, and I cannot generate all possible options to see the covered grammar.

A possible solution:
 - create small separated files and define LOCAL symbols for testing purpose
 - the small files are MERGED later, without the LOCAL sections, so the small tested units can be merged to a bigger one.

## Parser: Lark is selected
Lark can read BNF-like input, and produce an AST, so to save time/energy,
I use that now.


## Language grammar modification steps
 - change/create grammar_X_something.bnf files
   - use local variables to reduce the combinations of
     lower level language elems (otherwise you cannot generate
     a sample set)
 - what are the changes in grammar files?
    - use bnf_validator/bnf_validator to generate examples for the new rules to see what are the side effects of your changes
 - use erlango_lang_bnf_to_lark_merger.py to:
   - read all grammar_X_.bnf files
   - remove local sections (locals are defined in lower level grammar files)
   - convert them to 'erlango_lang_def.lark' file
   - generate a new Abstract Syntax Tree from the new lark file
