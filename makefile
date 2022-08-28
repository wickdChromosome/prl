build:
	go build .
	rm test/*.zip
	time ./prl -j 5 -cmd "zip -r {paths.txt}.zip {paths.txt}"
	#rm test/*.zip
	#time for i in `cat paths.txt`; do zip -r $i.zip $i; done

