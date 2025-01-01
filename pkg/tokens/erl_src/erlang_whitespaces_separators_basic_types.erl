#!/usr/bin/env escript

% test file for basic types,
% whitespaces, commas, dots, colons
% atom, string, integer, float, hexa

main(_) ->
	% tab used here as indentation
	io:fwrite("~p~n", [example()]),
	io:fwrite("~p~n", [example(1)]),
	io:fwrite("~p~n", [example(1234567890)]),

    io:fwrite("~p~n", [example(12.34)]),
    io:fwrite("~p~n", [add(2, 4)]),
    io:fwrite("~p~n", [double(9)]),

    io:fwrite("~p~n", [half(10)]),

    ok.

example() ->
    example.

% "double \"quoted\" comment, with a single quoted 'atom'"
example(1) -> "case 1 \\\" complex string \" with \n newline";                     % comment in 'example' function
example(1234567890) -> "case 1234567890";
example(16#af6bfa23) -> "hexa num";           % hexa based number
example(12.34) -> 12.34;           % hexa based number
example(atom_direct) -> atom_direct;
example('text_block_single_quoted') -> 'text_block_single_quoted_reply with \' escape and \n newline';
example(_) -> "case others".

add(A, B) ->
    Result = (A + B),
    Result.

diff(A, B) -> A - B.

double(A) -> A * 2.

half(B) -> B / 2.
