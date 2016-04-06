default: build

version := "v0.0.1"

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
name := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

ifeq ($(OS),Windows_NT)
	bin_suffix := ".exe"
else
	bin_suffix := ""
endif

clean:
	rm -f ./bin/$(name)*
	rm -f /usr/local/bin/$(name)*

compile:
	GOGC=off CGOENABLED=0 go build -i -o ./bin/$(name)$(bin_suffix) ./bin

print-success:
	@echo
	@echo "Plugin built."
	@echo
	@echo "To use it, either run 'make install' or set your PATH environment variable correctly."

build: compile print-success

release:
	rm -rf release
	mkdir release
	GOOS=linux GOARCH=amd64 GOGC=off CGOENABLED=0 go build -i -o release/$(name) ./bin
	tar --remove-files -cvzf release/$(name)-$(version)-linux-amd64.tar.gz -C release $(name)
	GOOS=darwin GOARCH=amd64 GOGC=off CGOENABLED=0 go build -i -o release/$(name) ./bin
	tar --remove-files -cvzf release/$(name)-$(version)-darwin-amd64.tar.gz -C release $(name)
	GOOS=windows GOARCH=amd64 GOGC=off CGOENABLED=0 go build -i -o release/$(name).exe ./bin
	tar --remove-files -cvzf release/$(name)-$(version)-windows-amd64.tar.gz -C release $(name).exe

install:
	cp bin/$(name) /usr/local/bin/


.PHONY : build release install