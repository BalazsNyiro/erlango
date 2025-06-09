# erlango

Erlang interpreter - written in Go
State machine replication - Distributed Virtual Machine,
Shared code execution on multiple nodes.


# Goals in the first session:
 - Erlang parser creation with basic language elems
 - Virtual Machine creation to execute the code
 - basic filesystem/string/network handling.

The first session doesn't want to be perfect.
It is a POC, and when it works, the code will be refactored/tuned.

# Documentation

The used Erlang documentation is v27.2
https://www.erlang.org/doc/system/data_types.html


# Install

## Lexer/Parser



### Parser: Lark
The grammar is defined in BNF files, so [Lark] (https://github.com/lark-parser/lark) was selected to generate the AST.
```
pip install lark==1.2.2 --upgrade
```


