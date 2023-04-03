#!/usr/bin/env escript

%INFO tested:
%INFO newline | tab | space | comma | dot | semicolon
       % Token_type_whitespace
          % Token_type_whitespace (newline as last char)
main(_) ->
% Token_type_whitespace
	io:format(example(), "tab indentation instead of space, before io:format").

                       %INFO Token_type_comma
                        % Token_type_whitespace (space)
    io:format(example(), "comma token in the line").
          % Token_type_whitespace (space)
             % Token_type_whitespace (space)
                      %INFO Token_type_semicolon
example(1) -> "case 1";
                           %INFO Token_type_dot
example(_) -> "case others".

