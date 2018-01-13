build:
	protoc -I. --go_out=plugins=micro:$(CURDIR) \
	proto/user/user.proto

	docker build -t user-service .

run:
	docker run -d --net="host" \
		-p 50051 \
		-e DB_HOST=localhost \
		-e DB_PASS=postgres \
		-e DB_USER=postgres \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns \
		user-service