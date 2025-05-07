PHONY: run-all  buf docker-push docker-build postgres createdb
run-all:
	cd auth-service && start /b make run-server
	cd order-service && start /b make run-server
	cd product-service && start /b  make run-server
	cd gateway && start /b make run-server
buf:
	@if exist pb\*.go del /Q pb\*.go
	buf generate
docker-build:
	docker build -f ./gateway/Dockerfile -t vietquandeptrai/api-gateway .
	docker build -f ./auth-service/Dockerfile -t  vietquandeptrai/auth-svc .
	docker build -f ./product-service/Dockerfile -t vietquandeptrai/product-svc .
	docker build -f ./order-service/Dockerfile -t vietquandeptrai/order-svc .

docker-push:
	cd gateway && docker push vietquandeptrai/api-gateway
	cd auth-service && docker push vietquandeptrai/auth-svc
	cd product-service && docker push vietquandeptrai/product-svc
	cd order-service && docker push vietquandeptrai/order-svc
postgres:
	@docker run --name auth-db -p 5431:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=12345 -d postgres
createdb:
	@docker exec -it auth-db createdb --username=postgres --owner=postgres auth_db


