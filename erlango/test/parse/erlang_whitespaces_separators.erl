#!/usr/bin/env escript

%INFO tested:
%INFO newline | tab | space | comma | dot | semicolon
       % Token_type_whitespace
          % Token_type_whitespace (newline as last char)
main(_) ->
% Token_type_whitespace  tab indented line:
	io:format(example()).

%INFO in the tab indented space it's difficult to represent the string, because the tab causes display problems
%INFO so I create a separated line to that
                         %%%%%%%%%%%%%%%%%%% Token_type_txt_quoted_double
    io:format(example(), "space indentation").

                       % Token_type_comma
                        % Token_type_whitespace (space)
    io:format(example(), "comma token in the line").
          % Token_type_whitespace (space)
             % Token_type_whitespace (space)
                      % Token_type_semicolon
example(1) -> "case 1";
                           % Token_type_dot
example(_) -> "case others".

