#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import unittest, bnf_lib, os


"""
Run all tests:
python3 bnf_lib_test.py 
"""

# python3 bnf_lib_test.py Test_symbol_class
class Test_symbol_class(unittest.TestCase):
    
    def test_symbol_obj_expand_possibilities(self):
        errors: list[str] = []
        filePathBnf = "../grammar_40_simple_types.bnf"

        symbolsTable, symbolNamesInLocalDefinition, errors = (
            bnf_lib.symbols_detect_in_file(filePathBnf, errors))

        wantedSymbolName = "<atomSmallFirstChar_tail>"
        self.assertIn(wantedSymbolName, symbolsTable)

        symbol = symbolsTable[wantedSymbolName]
        possibilities = symbol.expandPossibilities()
        bnf_lib.symbolnames_possibilities_print(possibilities, prefix="possibilities test:")

        # add an extra terminating, for testing reasons:
        extraTerminatings = ["t", "e", "r", "m", "i", "n", "a", "t", "i", "n", "g", "<empty>"]
        possibilities.append(extraTerminatings)
        self.assertEqual([['<atomPossibleCharAfterFirstPosition>', '<atomSmallFirstChar_tail>'], ['<empty>'], extraTerminatings],
                         possibilities)

        nonTermsFromAllPossib = symbol.grammar_elems_nonterminating_collect_in_all_possibilities()
        print(f"extra terminatings are NOT collected here")
        print(f"non terms: {nonTermsFromAllPossib}")
        self.assertEqual(['<atomPossibleCharAfterFirstPosition>', '<atomSmallFirstChar_tail>', '<empty>'],
                         nonTermsFromAllPossib)


# python3 bnf_lib_test.py Test_symbols_collect
class Test_symbols_collect(unittest.TestCase):

    def test_symbols_collect_from_grammar_basic(self):
        grammar = '<integer> "." <letter>'
        possibilities = bnf_lib.symbol_names_collect_from_grammar_def(grammar, verbose=True)
        symbols = possibilities[0]
        print(f"symbols: {symbols}")

        # python3 bnf_lib_test.py Test_symbols
        self.assertTrue(len(symbols) == 3)
        
        wanted = ['<integer>', '"."', '<letter>']
        self.assertEqual(wanted, symbols)


    # python3 bnf_lib_test.py Test_symbols_collect.test_symbols_collect_from_grammar_escape
    def test_symbols_collect_from_grammar_escape(self):
        grammar = f'"\\"" <letter> "\\""'
        print(f"GRAMMAR, escapeTest: {grammar}")

        possibilities = bnf_lib.symbol_names_collect_from_grammar_def(grammar)
        symbols = possibilities[0]
        print(f"symbols: {symbols}")

        # python3 bnf_lib_test.py Test_symbols
        self.assertTrue(len(symbols) == 3)

        wanted = ['"\\""', '<letter>', '"\\""']
        self.assertEqual(wanted, symbols)


    # python3 bnf_lib_test.py Test_symbols_collect.test_symbols_collect_nonterminatings_only
    def test_symbols_collect_nonterminatings_only(self):
        grammar = '<integer> "." <integer> "\\"" <letter> "\\""  "end"'

        possibilities= bnf_lib.symbol_names_collect_from_grammar_def(grammar)
        symbols = possibilities[0]
        symbolsNonTerm = bnf_lib.symbols_nonterminating_collect(symbols)

        wanted = ['<integer>', '"."', '<integer>',  '"\\""', '<letter>', '"\\""', '"end"']
        self.assertEqual(wanted, symbols)

        wanted = ['<integer>', '<integer>', '<letter>']
        self.assertEqual(wanted, symbolsNonTerm)





# python3 bnf_lib_test.py Test_symbols
class Test_symbols(unittest.TestCase):

    def test_symbolname_concatenate(self):
        """only for cover, this is a simple one-liner"""
        self.assertEqual("a, b", bnf_lib.symbolnames_concate_simple_str(["a", "b"]))


    # python3 bnf_lib_test.py Test_symbols.test_symbolname_grammar_definition_in_file
    def test_symbolname_grammar_definition_in_file(self):
        """collect new symbolname and grammar from lines"""

        errors: list[str] = []

        line =    '<symbol> ::= "a" | <other>'
        defOnly = '             "a" | <other>'
        # the definition's indentation is kept, because if the grammar has multiple lines,
        # the multi-line display can keep the formatting then

        newSymbolNameInLine, definitionInLine, errors = bnf_lib.symbolname_and_grammar_definition_in_line__get(line, errors)

        self.assertEqual(newSymbolNameInLine, "<symbol>")
        self.assertEqual(definitionInLine, defOnly)


        line = 'invalid-symbol-missing-opening-bracket> ::= "a" | <other>'
        newSymbolNameInLine, definitionInLine, errors = bnf_lib.symbolname_and_grammar_definition_in_line__get(line, errors)
        # the last error is...
        self.assertTrue(errors[-1].startswith("missing symbol brackets"))


        line = '<%%% incorrect character in symbolname!!> ::= "a" | <other>'
        newSymbolNameInLine, definitionInLine, errors = bnf_lib.symbolname_and_grammar_definition_in_line__get(line, errors)
        # the last error is...
        self.assertTrue(errors[-1].startswith("maybe human error"))


    def test_is_terminating_symbolname(self):

        name = '<atom>'
        self.assertFalse(bnf_lib.symbolname_terminating(name))

        name = '"atom>'
        self.assertFalse(bnf_lib.symbolname_terminating(name))

        name = '"a"'
        self.assertTrue(bnf_lib.symbolname_terminating(name))


    def test_symbols_detect_in_file(self):
        filePathBnf = "../grammar_40_simple_types.bnf"
        print(f"filePath: {filePathBnf}")
        errors: list[str] = []

        symbols, symbolNamesInLocalDefinition, errors = (
            bnf_lib.symbols_detect_in_file(filePathBnf, errors))

        symbolNamesInFile = sorted(symbols.keys())
        print(f"detected symbolNames: {symbolNamesInFile}")
        symbolNamesWanteds = [
            '<anyUnicodeCharExceptDoubleQuote>',
            '<anyUnicodeCharsExceptSingleQuote>',

            '<atom>', '<atomCharList_inQuotes>', '<atomInQuotes>', '<atomPossibleCharAfterFirstPosition>',
            '<atomSmallFirstChar>', '<atomSmallFirstChar_tail>',

            '<digit>',
            '<empty>',
            '<float>',
            '<integer>',
            '<letter>',
            '<letterSmall>',
            '<number>',
            '<numberTail>',
            '<pid>',

            '<triple_anyCharExceptQuoteOrEmtpy>', '<stringQuoteTriple_safeChar>',
            '<string>', '<stringQuoteTriple>', '<stringQuoteOne>', '<stringQuoteOneTail>', '<stringQuoteTripleTail>'
        ]

        self.assertEqual(symbolNamesWanteds, symbolNamesInFile)
        bnf_lib.symbols_table_print(symbols)

        self.assertIn("<letter>", symbolNamesInLocalDefinition)
        self.assertNotIn("<atom>", symbolNamesInLocalDefinition)



    def test_symbols_double_definition_detect(self):
        filePathBnf = "bnf_lib_test_grammar_double_symbol.bnf"
        print(f"filePath: {filePathBnf}")
        errors: list[str] = []

        symbols, symbolNamesInLocalDefinition, errors = (
            bnf_lib.symbols_detect_in_file(filePathBnf, errors))

        print(f"errors: {errors}")
        
        self.assertEqual(errors, ['ERROR: the symbol is defined more than once in the bnf grammar: <letter>, defCount: 2 '])
        self.assertEqual(symbols["<letter>"].definitionCounterInBnf, 2)


class Test_file_related_funcs(unittest.TestCase):
    print("""give cover to file related funcs""")

    # for coverage
    def test_file_read_write(self):

        fname = "/tmp/test_file_read_write.txt"
        content = "abcdefgh"
        bnf_lib.file_write(fname, content)
        print(f"file written: {fname} <- {content}")

        collectedFiles = bnf_lib.files_collect_in_dir("/tmp", prefix="test_file")
        print(f"collected files: {collectedFiles}")

        for fileCollected in collectedFiles:
            print(f"test, file collected: {fileCollected}")

            readBack = []
            if bnf_lib.file_exists___alert_if_not(fileCollected):
                readBack = bnf_lib.file_src_lines(fileCollected)
                print(f"readBack: {readBack}")
            self.assertEqual([content], readBack)

        os.remove(fname)
        self.assertEqual(False, bnf_lib.file_exists___alert_if_not(fname, raiseException=False))
        self.assertEqual(False, bnf_lib.file_exists___alert_if_not(fname, raiseException=False))
        with self.assertRaises(ValueError):
            # exception is raised by default
            bnf_lib.file_exists___alert_if_not(fname)


if __name__ == '__main__':  # pragma: no cover
    unittest.main()
