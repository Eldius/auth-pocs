
basic:
	$(MAKE) -C helper-library build
	$(MAKE) -C basic-auth start

build-docker:
	$(MAKE) -C basic-auth build-docker
