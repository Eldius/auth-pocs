
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
		go test -C "$$dir" ./... -cover || exit 1; \
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
