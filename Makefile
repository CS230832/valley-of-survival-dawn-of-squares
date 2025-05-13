run: build
	docker run -d -p 8080:8080 --name vosdos vosdos

build: clean
	docker build -t vosdos .

clean:
	-@docker stop vosdos 2>/dev/null || true
	-@docker rm vosdos 2>/dev/null || true
	-@docker rmi vosdos 2>/dev/null || true