
basic:
	$(MAKE) -C helper-library build
	$(MAKE) -C basic-auth start

build-docker:
	$(MAKE) -C basic-auth build-docker

tidy:
	$(MAKE) -C helper-library tidy
	$(MAKE) -C basic-auth tidy