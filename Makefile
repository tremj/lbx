test:
	docker-compose -f docker-compose.test.yaml up -d
	go test ./...
	docker-compose -f docker-compose.test.yaml down