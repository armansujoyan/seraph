COMPILER = ./seraph

all:
	go build -o seraph src/main.go

test:
	$(COMPILER) test.pas
	./test
	echo $?
