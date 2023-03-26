#!/usr/bin/env escript

% tested:
% newline | tab | space | comma | dot | semicolon
       % Token_type_whitespace
          % Token_type_whitespace (newline as last char)
main(_) ->
% Token_type_whitespace (tab indentation)
                                  % Token_type_comma
	io:format(  example_fun_call(), "something else").

                   % Token_type_whitespace (space)
                      % Token_type_whitespace (space)
                               % Token_type_semicolon
example_fun_call(1) -> "case 1";
                                    % Token_type_dot
example_fun_call(_) -> "case others".

