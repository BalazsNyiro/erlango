-module(hello).
-export([hello/0]).

hello() ->
  io:format("Direct printing \" quote from Erlang code ~n"),
  io:format(example_fun_call()).

example_fun_call() ->
  Msg = "example fun call output~n",
  Msg ++ "~nmsg end.~n".



