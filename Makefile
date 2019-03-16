.PHONY : all

all:
	build/create_environment.sh

clean:
	rm -rf build/go_workspace
	rm -rf eth/build
	rm -rf eth/gen
	rm filetribe