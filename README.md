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

Documentation

# Compile this project

Remove older go version:
```
apt remove gccgo-go
```

Install Go: 
https://go.dev/wiki/Ubuntu
```
sudo snap install --classic go
```

Test it in a new terminal
```
go version
go version go1.23.4 linux/amd64
```

