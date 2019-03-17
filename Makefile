.PHONY : all

all:
	build/install.sh

clean:
	rm -rf build/go_workspace
	rm -rf eth/build
	rm -rf eth/gen
	rm build/bin/filetribe