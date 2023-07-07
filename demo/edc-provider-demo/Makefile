

provider-up:
	docker-compose -f ./organizations/provider/provider.docker-compose.yaml  up
provider-down:
	docker-compose -f ./organizations/provider/provider.docker-compose.yaml  down

consumer-up:
	docker-compose -f ./organizations/consumer/consumer.docker-compose.yaml  up
consumer-down:
	docker-compose -f ./organizations/consumer/consumer.docker-compose.yaml  down

demo-up:
	provider-up
	consumer-up

demo-down:
	provider-down
	consumer-down
