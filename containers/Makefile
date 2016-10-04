.PHONY: docker-build

build:
	make -C hellogo build
	make -C hellonode build
	make -C simplerpc build
	make -C auth build

build-docker:
	make -C hellogo build-docker
	make -C hellonode build-docker
	make -C simplerpc build-docker
	make -C auth build-docker
