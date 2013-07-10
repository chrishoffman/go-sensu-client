#!/bin/sh

mockgen -package=mocks sensu-client/sensu MessageQueuer >sensu/mocks/mock_message_queuer.go
