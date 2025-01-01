#!/usr/bin/env escript

-mode(compile).
% escript must be compiled to use fun in foldl

main(_unused_var_starts_with_underscore) ->

  Text = <<"progress 99%: msg">>, % comment after a "string"
  % comment in an empty line
  
  io:format(Text),

  MultilineText = "
                   This is a string,
                   that contains newlines
                   between quotes",

  io:format("multi-line text: ~n"),
  io:format(MultilineText),


  ThisIsNotAcommentAtEndOfLine= "\n\nab % cde
              fgh\n",
  io:format(ThisIsNotAcommentAtEndOfLine),

  io:format(example_fun_call()).

example_fun_call() ->
  "example fun call output~n".

