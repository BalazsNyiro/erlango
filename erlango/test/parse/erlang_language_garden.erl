#!/usr/bin/env escript

-mode(compile).
% escript must be compiled to use fun in foldl

main(_) ->

  Text = <<"progress 99%: msg">>, % comment after a "string"
  % comment in an empty line
  
  io:format(Text),

  MultilineText = "
                   This is a string,
                   that contains newlines
                   between quotes",


  io:format("multi-line text: ~n"),
  io:format(MultilineText),


  Comments = "\n\nab % cde
              fgh",
  io:format(Comments),

  io:format(example_fun_call()).

example_fun_call() ->
  "example fun call output~n".

