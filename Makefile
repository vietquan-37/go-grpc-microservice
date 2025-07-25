PHONY: run-all  buf docker-push docker-build auth-db product-db order-db
run-all:
	cd auth-service && make run-server &
	cd product-service && make run-server &
	cd order-service && make run-server &
	cd payment-service && make run-server &
	cd email-service && make run-server &
	cd gateway && make run-server &
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
auth-db:
	@docker run --name auth-db -p 5431:5432 \
		-e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=12345 \
		-e POSTGRES_DB=auth_db \
		-d postgres:16.4

product-db:
	@docker run --name product-db -p 5430:5432 \
		-e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=12345 \
		-e POSTGRES_DB=product_db \
		-d postgres:16.4

order-db:
	@docker run --name order-db -p 5429:5432 \
		-e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=12345 \
		-e POSTGRES_DB=order_db \
		-d postgres:16.4
kill-go:
	@pkill -f 'go run' || true
