.PHONY: \
run help

.DEFAULT_GOAL := help

run: ## 시작
	go run ./cmd/conflunce

help: ## 항목
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'