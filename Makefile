# I personally use podman, but I added docker support just in case.
PODMAN_VERSION_MAJOR=$(shell podman -v | sed -r 's/^podman version ([0-9]+)\..*$$/\1/')
ifeq ($(shell expr $(PODMAN_VERSION_MAJOR) \>= 4), 1)
    CONTAINER_RUNTIME=podman
else
    $(warning Podman found but too old ($(PODMAN_VERSION_MAJOR)). Falling back to docker.)
    CONTAINER_RUNTIME=docker
endif

.PHONY: up down

up:
	$(CONTAINER_RUNTIME) compose -f docker-compose.yml up -d

down:
	$(CONTAINER_RUNTIME) compose -f docker-compose.yml down

run:
	go run . serve

