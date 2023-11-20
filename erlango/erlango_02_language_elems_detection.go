/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.

Version 0.2, second rewrite

*/

package erlango




/*

== Definitions for Erlang language elem detection ==

https://www.erlang.org/doc/reference_manual/expressions.html

Erlang term (data types):
	- an integer, float, atom, string, list, map, or tuple.

Erlang variables
	Variables start with an uppercase letter or underscore (_). Variables can contain alphanumeric characters, underscore and @.

	Variables starting with underscore (_), for example, _Height, are normal variables, not anonymous.

	Special chars are NOT allowed:
	Eshell V13.1.5  (abort with ^G)
	1> Aáéői = 3.
	* 1:4: illegal character

Operators:
	whitespaces are not important in operator detection:
	1> A = 2+-1.
	1




== TOKENS ==
 - meaningful characters that can be interpreted in different environments.
   Tokens have to be interpreted with their environments/positions
   For example: a . have a different meanings in these situations
		- 1.2
		- #Name.Field
		- . at the end of a function.

Language elem detection steps:

 - detect simple terms:
    - integers: https://www.erlang.org/doc/reference_manual/data_types.html
		- 1234.
		- 1_234_567_890.
		- $A.
		- 16#1f.
		- 16#4865_316F_774F_6C64.
		- 2e-3.

	- floats:
		- 2.3.
		- 2.3e-3.
		- 1_234.333_333.

	- atom-quoted
	- atom

	- string

 - detect more complex elems

	- bit-string https://www.erlang.org/doc/man/binary.html
	-


 - detect complex terms: list, map, tuple,





== LANGUAGE ELEMS


== LINKS, REFERENCES ==
           data types:  https://www.erlang.org/doc/reference_manual/data_types
dot, colon, semicolon:  https://stackoverflow.com/questions/1110601/in-erlang-when-do-i-use-or-or
   terms, expressions:  https://www.erlang.org/doc/reference_manual/expressions#terms


== NUMBERS ==
1> 1_000 * 3.
3000




*/