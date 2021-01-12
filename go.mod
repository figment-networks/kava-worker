module github.com/figment-networks/kava-worker

go 1.15

replace github.com/figment-networks/indexer-manager => /Users/pacmessica/.go/src/github.com/figment-networks/indexer-manager

require (
	github.com/bearcherian/rollzap v1.0.2
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/figment-networks/indexer-manager v0.0.8
	github.com/figment-networks/indexing-engine v0.1.14
	github.com/google/uuid v1.1.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/rollbar/rollbar-go v1.2.0
	github.com/sirupsen/logrus v1.7.0
	go.uber.org/zap v1.16.0
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324
	google.golang.org/grpc v1.34.0
)
