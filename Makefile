.SILENT:

run:
	go run ./cmd/main.go

docBuild:
	docker build -t forum .

docrun: docBuild
	docker run -it -p 8080:8080 forum

docDelete:
	docker rmi forum

docClear: docDelete
	docker system prune -a