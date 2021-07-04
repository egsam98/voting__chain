module github.com/egsam98/voting/chain

go 1.16

require (
	github.com/Shopify/sarama v1.29.0
	github.com/egsam98/voting/proto v0.2.0
	github.com/go-chi/chi/v5 v5.0.3
	github.com/go-chi/render v1.0.1
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/joho/godotenv v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.23.0
	google.golang.org/protobuf v1.26.0
)

replace github.com/egsam98/voting/proto => github.com/egsam98/voting__proto v0.2.0
