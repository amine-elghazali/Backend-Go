# Using a Makefile to simplify and automate certain tasks in your development process.
# Here we're using make build , to run the ' go build -o bin/backend_go ' via this simple cli command : make build  .  then what's happened, it's runs the following command first : go build -o bin/backend_go  that we have in the build, than execute the command we have in run, which is just executing the following go binary generated file   :  ./bin/backend_go
build : 
	go build -o bin/backend_go
run : build 
	./bin/backend_go
test : 
	go test -v ./...
# tests in the current directory and its subdirectories 