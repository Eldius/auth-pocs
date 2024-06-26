
MODULES := helper-library basic-auth

basic:
	$(MAKE) -C helper-library build
	$(MAKE) -C basic-auth start

build-docker:
	$(MAKE) -C basic-auth build-docker

tidy:
	$(MAKE) -C helper-library tidy
	$(MAKE) -C basic-auth tidy

vulncheck:
	for dir in $(MODULES); do \
	    echo "#####################"; \
	    echo "# starting for $$dir"; \
	    echo "#####################"; \
		govulncheck -C "$$dir" ./... || exit 1; \
	    echo ""; \
	    echo "ending for $$dir"; \
	    echo ""; \
	    echo "---------------------"; \
	done

lint:
	$(eval CURR_PATH := $(PWD))
	for dir in $(MODULES); do \
	    echo "#####################"; \
	    echo "# starting for $$dir"; \
	    echo "#####################"; \
		cd "$(CURR_PATH)/$$dir"; golangci-lint run || exit 1; \
	    echo ""; \
	    echo "ending for $$dir"; \
	    echo ""; \
	    echo "---------------------"; \
	done
	cd $(CURR_PATH)


test:
	for dir in $(MODULES); do \
	    echo "#####################"; \
	    echo "# starting for $$dir"; \
	    echo "#####################"; \
		go test -C "$$dir" -coverprofile=coverage.out ./... -cover || exit 1; \
		go tool -C "$$dir" cover -html=coverage.out || exit 1; \
	    echo ""; \
	    echo "ending for $$dir"; \
	    echo ""; \
	    echo "---------------------"; \
	done

root:
	curl -i -u "root:12345" localhost:8080/api/endpoint

user:
	curl -i -u "user:12345" localhost:8080/api/endpoint

invalid-pass:
	curl -i -u "root:1234" localhost:8080/api/endpoint

invalid-user:
	curl -i -u "invalid:12345" localhost:8080/api/endpoint

no-auth:
	curl -i localhost:8080/api/endpoint

basic-docker-env-down:
	docker compose -f docker-compose-basic.yml down

basic-docker-run: basic-docker-env-down
	LOG_LEVEL=debug docker compose -f docker-compose-basic.yml up --build --remove-orphans -d

query-lok-logs:
	curl -G -s  "http://localhost:3100/loki/api/v1/query" \
      --data-urlencode 'query=sum(rate({job="varlogs"}[10m])) by (level)' | jq
