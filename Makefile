install:
	go build .
	rm ..\..\Desktop\developer-tools.exe
	mv .\developer-tools.exe '..\..\Desktop\developer-tools.exe'