run:
	go build -v && DEPTH=3 ROOT=http://0.0.0.0:8080/ ./scrapper && rm scrapper

site:
	cd test-site && ./run