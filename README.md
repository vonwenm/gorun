# Introdution

Auto build and execute code when detect file changing.

# Example

### Normal usage

Change to you project path and execute command "gorun"

	cd $GOPATH/src/myProject
	gorun

It will be

	go build
	./myProject

### With arguments

	gorun -p /tmp -f -c cat /tmp/access.log

It will be

	go build
	./myProject -p /tmp -f -c cat /tmp/access.log


