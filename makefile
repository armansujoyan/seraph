COMPILER = ./seraph

all:
	go build -o seraph src/main.go

clean:
	rm $(COMPILER) *.s

test:
	$(COMPILER) test.pas
	./test
	echo $?
