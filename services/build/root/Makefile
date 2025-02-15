include build/root/config.mk
include .env

BINS ?= auth
IMAGE_DOTFILES = $(foreach bin,$(BINS),.image-$(bin)-$(TAG))

BUILD_DIRS := .dist/$(OS)_$(ARCH)	\
			  .images				\
			  .dockerfile			\
              .go/cache            	\
              .go/pkg

all: # @HELP builds binaries for one platform ($OS/$ARCH)
all: build images push

build: # @HELP for each BIN, build binary and place into /bin/$(OS)_$(ARCH)
build: $(BINS)

# Build each binary using Docker
$(BINS): $(BUILD_DIRS)
	docker run                                                  \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
		-v $$(pwd)/.dist/$(OS)_$(ARCH):/.dist                   \
		-v $$(pwd)/.go/cache:/.cache							\
	    --env GOCACHE="/.cache/gocache"                         \
	    --env GOMODCACHE="/.cache/gomodcache"                   \
	    --env ARCH="$(ARCH)"                                    \
	    --env OS="$(OS)"                                        \
	    $(CONTAINER_IMAGE)                                      \
	    go build -o .dist/$(OS)_$(ARCH)/"$@" ./cmd/"$@"

$(BUILD_DIRS):
	mkdir -p $@

clean: # @HELP clean artifacts from build
clean:
	rm -rf .go .dist .images .dockerfile

image images: # Builds container images
image images: $(IMAGE_DOTFILES)
	for bin in $(BINS); do \
	    echo "image: $(REGISTRY)/$$bin:$(TAG)";  \
	done
	echo

# Pattern rule to create .image files based on each BIN
.image-%-$(TAG): BIN = $(patsubst .image-%-$(TAG),%,$@)
.image-%-$(TAG):
	echo $(BIN)
	sed                                            \
		-e 's|{ARG_BIN}|$(BIN)$(BIN_EXTENSION)|g'  \
		-e 's|{ARG_ARCH}|$(ARCH)|g'                \
		-e 's|{ARG_OS}|$(OS)|g'                    \
		-e 's|{ARG_FROM}|$(CONTAINER_IMAGE)|g'     \
		./build/Dockerfile.in > ./.dockerfile/.dockerfile-$(BIN)-$(OS)_$(ARCH)
	docker build \
		--progress=plain \
		-t $(REGISTRY)/$(BIN):$(TAG) \
	    -f ./.dockerfile/.dockerfile-$(BIN)-$(OS)_$(ARCH) \
		.
	docker images -q $(REGISTRY)/$(BIN):$(TAG) > ./.images/$@
	echo

push: # @HELP pushes the image for one platform ($OS/$ARCH) to the defined registry
push: image
	for bin in $(BINS); do                     \
	    docker push $(REGISTRY)/$$bin:$(TAG);  \
	done
	echo

lint: # HELP runs golangci-lint linters
	@if golangci-lint run; then \
		echo "golangci-lint ran successfully"; \
	else \
		echo "golangci-lint failed to lint project files"; \
		read -p "Would you like to try running it again with --fix (y/n)? " selection; \
		if [ "$$selection" = "y" ]; then \
			echo "Running golangci-lint with --fix..."; \
			golangci-lint run --fix; \
			if [ $$? -eq 0 ]; then \
				echo "golangci-lint ran successfully with --fix"; \
			else \
				echo "golangci-lint still failed to lint project files."; \
			fi; \
		else \
			echo "Skipping fix attempt. Please review the linting errors."; \
		fi; \
	fi

run: # HELP select and run a single application
run:
	@echo "$(BINS)"
	@read -p "Select an application to run: " selection; \
	if echo "$(BINS)" | grep -wq $$selection; then \
	    $(PWD)/build/scripts/run-locally.sh $$selection \
	else \
	  printf  "application: \"%s\", does not exist.\n" $$selection; \
	fi

####################################################################################################
# Database Tasks
# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#migrate-cli
####################################################################################################

db/connect: # @HELP connect to locally running database
db/connect:
	$(MAKE) -C ./migrations db/connect/psql

db/migrations/up: # @HELP run migrations to create database tables
db/migrations/up:
	$(MAKE) -C ./migrations db/migrations/up

db/migrations/down: # @HELP run migrations to teardown database tables
db/migrations/down:
	$(MAKE) -C ./migrations db/migrations/down

db/migrations/refresh: # @HELP refresh database tables
db/migrations/refresh: db/migrations/down db/migrations/up

SHELL := /usr/bin/env bash -o errexit -o pipefail
.DEFAULT_GOAL = all
.PHONY: all clean