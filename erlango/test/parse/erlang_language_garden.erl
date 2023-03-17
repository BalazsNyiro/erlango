#!/usr/bin/env escript

-mode(compile).
% escript must be compiled to use fun in foldl

main(_) ->

  Text = <<"progress 99%: msg">>, % comment after a "string"
  % comment in an empty line

  io:format(Text),
  io:format(" - Direct printing from Erlang code ~n"),
  io:format(example_fun_call()).

example_fun_call() ->
  "example fun call output~n".

