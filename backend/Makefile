GO_FILES=$(shell go list ./... | grep -v 'mocks' | grep -v '^.*fake.*' | grep -v '/grpc')
MOCKGEN=go run github.com/golang/mock/mockgen
DELIEVERY_PATH=internal/delivery
USECASE_PATH=internal/usecase
REPOSITORY_PATH=internal/repository

.PHONY: coverage
coverage:
	go test $(GO_FILES) -coverprofile=coverage.out

.PHONY: summarize-coverage
summarize-coverage: coverage
	go tool cover -func=coverage.out

.PHONY: mockgen
mockgen:
	${MOCKGEN} -source=$(DELIEVERY_PATH)/handlers/auth-handler.go -destination=$(DELIEVERY_PATH)/mocks/auth-mock.go -package=mocks
	${MOCKGEN} -source=$(DELIEVERY_PATH)/handlers/pvz-handler.go -destination=$(DELIEVERY_PATH)/mocks/pvz-mock.go -package=mocks
	${MOCKGEN} -source=$(DELIEVERY_PATH)/handlers/reception.go -destination=$(DELIEVERY_PATH)/mocks/reception-mock.go -package=mocks

	${MOCKGEN} -source=$(USECASE_PATH)/auth-usecase.go -destination=$(USECASE_PATH)/mocks/auth-mock.go -package=mocks
	${MOCKGEN} -source=$(USECASE_PATH)/pvz-usecase.go -destination=$(USECASE_PATH)/mocks/pvz-mock.go -package=mocks
	${MOCKGEN} -source=$(USECASE_PATH)/reception-usecase.go -destination=$(USECASE_PATH)/mocks/reception-mock.go -package=mocks

.PHONY: integration
integration:
	go test -tags=integration ./...

.PHONY: gen-proto
gen-proto:
	cd internal/grpc/pvz && protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto
