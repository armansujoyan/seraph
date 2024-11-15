COMPILER = ./compiler

all:
	go build src/main.go

test:
	$(COMPILER) test.pas
	./test
	echo $?
