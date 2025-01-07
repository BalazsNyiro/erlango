echo "cover create in current directory (run this where you have go files)"
pwd
go test -coverprofile=cover.out; go tool cover -html=cover.out

