include .env
export $(shell sed 's/=.*//' .env)

GOPATH=$(shell go env GOPATH)

test:
	@ echo
	@ echo "Testing Helm..."
	@ echo
	@ go run cmd/main.go $(filter-out $@,$(MAKECMDGOALS))
