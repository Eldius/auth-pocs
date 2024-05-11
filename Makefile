
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
