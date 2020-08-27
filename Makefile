ANT=antlr4

antlr: parser parser/zbxcsv
	$(ANT) -visitor -Dlanguage=Go -package zbxcsv -o parser/zbxcsv zbxcsv.g4


test:
	go test -run TestBasecommandsVisitor