module Payments

go 1.25.0

require (
	github.com/erooohaaa/orders-generated v0.0.0-20260412110537-4071211e7c4b
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.12.1
	github.com/rabbitmq/amqp091-go v1.10.0
	google.golang.org/grpc v1.80.0
)

require (
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/erooohaaa/orders-generated => ../orders-generated