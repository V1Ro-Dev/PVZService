GO_FILES=$(shell go list ./... | grep -v 'mocks')
MOCKGEN=go run github.com/golang/mock/mockgen
DELIEVERY_PATH=internal/delivery
USECASE_PATH=internal/usecase
REPOSITORY_PATH=internal/repository


coverage.out:
	go test $(GO_FILES) -coverprofile=coverage.out

.PHONY: summarize-coverage
summarize-coverage: coverage.out
	go tool cover -func=coverage.out

.PHONY: mockgen
mockgen:
	${MOCKGEN} -source=$(DELIEVERY_PATH)/handlers/auth-handler.go -destination=$(DELIEVERY_PATH)/http/mocks/auth-mock.go -package=mocks
	${MOCKGEN} -source=$(DELIEVERY_PATH)/handlers/pvz-handler.go -destination=$(DELIEVERY_PATH)/http/mocks/pvz-mock.go -package=mocks
	${MOCKGEN} -source=$(DELIEVERY_PATH)/handlers/reception.go -destination=$(DELIEVERY_PATH)/http/mocks/reception-mock.go -package=mocks

	${MOCKGEN} -source=$(USECASE_PATH)/auth-usecase.go -destination=$(USECASE_PATH)/mocks/auth-mock.go -package=mocks
	${MOCKGEN} -source=$(USECASE_PATH)/pvz-usecase.go -destination=$(USECASE_PATH)/mocks/pvz-mock.go -package=mocks
	${MOCKGEN} -source=$(USECASE_PATH)/reception-usecase.go -destination=$(USECASE_PATH)/mocks/reception-mock.go -package=mocks

	${MOCKGEN} -source=$(REPOSITORY_PATH)/auth-repository.go -destination=$(REPOSITORY_PATH)/mocks/auth-mock.go -package=mocks
	${MOCKGEN} -source=$(REPOSITORY_PATH)/pvz-repository.go -destination=$(REPOSITORY_PATH)/mocks/pvz-mock.go -package=mocks
	${MOCKGEN} -source=$(REPOSITORY_PATH)/reception-repository.go -destination=$(REPOSITORY_PATH)/mocks/reception-mock.go -package=mocks
