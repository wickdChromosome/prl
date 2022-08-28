build:
	go build .

test:
	rm test/*.zip
	time ./prl -j 5 -cmd "zip -r {paths.txt}.zip {paths.txt}"
install:
	sudo cp prl /usr/local/bin/prl
