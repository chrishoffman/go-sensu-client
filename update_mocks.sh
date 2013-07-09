#!/bin/sh

mockgen -package=mocks sensu-client/sensu MessagePublisher,MessageConsumer >sensu/mocks/mock_rabbitmq.go
