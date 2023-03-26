#!/usr/bin/env escript

-mode(compile).
% escript must be compiled to use fun in foldl

% == Cascade Token Type Definition: ==

% the percent signs select character ranges
% 2nd visible section: Token_type_.......  (prefix is always fix!  token type definition)

% if the 2nd visible elem pre
% if the line visible chars start with a % and a token type is the next elem,
% then that line is a Cascade Token Type definition for more chars.

% after the second visible elem you can insert any comment into the line

% below the percent patterns the first line which doesn't start with %, is a valid line

%%%% Token_type_atom  <- a function name is a token - and this is a comment in a Cascade Token Type Definition line
    % Token_type_(
     %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% Token_type_atom
                                       % Token_type_)
                                        % Token_type_whitespace
                                         %% Token_type_fun_open
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


  Comments = "\n\nab % cde
              fgh",
  io:format(Comments),

  io:format(example_fun_call()).

example_fun_call() ->
  "example fun call output~n".

