#!/bin/sh

mockgen -package=mocks sensu-client/sensu Publisher >sensu/mocks/mock_publisher.go
