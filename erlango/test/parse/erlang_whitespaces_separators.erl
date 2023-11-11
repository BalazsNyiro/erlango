#!/usr/bin/env escript

main(_) ->
	% tab used here as indentation
	io:fwrite("~p~n", [example()]),
	io:fwrite("~p~n", [example(1)]),
	io:fwrite("~p~n", [example(1234567890)]),

    io:fwrite("~p~n", [add(2, 4)]),

    ok.

example() ->
    example.

% "double \"quoted\" comment, with a single quoted 'atom'"
example(1) -> "case 1";                     % comment in 'example' function
example(1234567890) -> "case 1234567890";
example(16#af6bfa23) -> "hexa num";           % hexa based number
example(12.34) -> "floated num";           % hexa based number
example(_) -> "case others".

add(A, B) ->
    Result = (A + B),
    Result.

diff(A, B) ->
    A - B.