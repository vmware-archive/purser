travis-build: install-plugin install-controller

install-plugin:
	go install github.com/vmware/purser/cmd/purser_plugin

install-controller: build

# The binary to build (just the basename).
BIN := purser_controller

# This repo's root import path (under GOPATH)
PRO := github.com/vmware/purser
DEP := vendor
PKG := pkg/purser_controller
CMD := cmd/purser_controller

# Where to push the docker image.
REGISTRY ?= gurusreekanth

# Which architecture to build - see $(ALL_ARCH) for options.
ARCH ?= amd64

# This version-strategy uses git tags to set the version string
#VERSION := $(shell git describe --tags --always --dirty)
#
# This version-strategy uses a manual value to set the version string
VERSION := 1.0.0

###
### These variables should not need tweaking.
###

ALL_ARCH := amd64 arm arm64 ppc64le

# Set default base image dynamically for each arch
ifeq ($(ARCH),amd64)
    BASEIMAGE?=alpine
endif
ifeq ($(ARCH),arm)
    BASEIMAGE?=armel/busybox
endif
ifeq ($(ARCH),arm64)
    BASEIMAGE?=aarch64/busybox
endif
ifeq ($(ARCH),ppc64le)
    BASEIMAGE?=ppc64le/busybox
endif

IMAGE := $(REGISTRY)/$(BIN)-$(ARCH)

BUILD_IMAGE ?= golang:1.8-alpine

DOCKER_MOUNT_MODE=delegated

# If you want to build all binaries, see the 'all-build' rule.
# If you want to build all containers, see the 'all-container' rule.
# If you want to build AND push all containers, see the 'all-push' rule.
all: build

build-%:
	@$(MAKE) --no-print-directory ARCH=$* build

container-%:
	@$(MAKE) --no-print-directory ARCH=$* container

push-%:
	@$(MAKE) --no-print-directory ARCH=$* push

all-build: $(addprefix build-, $(ALL_ARCH))

all-container: $(addprefix container-, $(ALL_ARCH))

all-push: $(addprefix push-, $(ALL_ARCH))

build: bin/$(ARCH)/$(BIN)

bin/$(ARCH)/$(BIN): build-dirs
	@echo "building: $@"
	@docker run                                                            \
	    -ti                                                                \
	    -u $$(id -u):$$(id -g)                                             \
	    -v $$(pwd)/.go:/go:$(DOCKER_MOUNT_MODE)                            \
	    -v $$(pwd)/$(CMD):/go/src/$(PRO)/$(CMD):$(DOCKER_MOUNT_MODE)                     \
	    -v $$(pwd)/$(PKG):/go/src/$(PRO)/$(PKG):$(DOCKER_MOUNT_MODE)                     \
	    -v $$(pwd)/$(DEP):/go/src/$(PRO)/$(DEP):$(DOCKER_MOUNT_MODE)                     \
	    -v $$(pwd)/bin/$(ARCH):/go/bin:$(DOCKER_MOUNT_MODE)                \
	    -v $$(pwd)/bin/$(ARCH):/go/bin/linux_$(ARCH):$(DOCKER_MOUNT_MODE)  \
	    -v $$(pwd)/.go/std/$(ARCH):/usr/local/go/pkg/linux_$(ARCH)_static:$(DOCKER_MOUNT_MODE)  \
	    -w /go/src                                                 \
	    $(BUILD_IMAGE)                                                     \
	    /bin/sh -c "                                                       \
	        ARCH=$(ARCH)                                                   \
	        VERSION=$(VERSION)                                             \
	        PKG=$(PKG)                                                     \
	        ./$(PRO)/$(CMD)/build/build.sh                                               \
	    "

DOTFILE_IMAGE = $(subst :,_,$(subst /,_,$(IMAGE))-$(VERSION))

container: .container-$(DOTFILE_IMAGE) container-name
.container-$(DOTFILE_IMAGE): bin/$(ARCH)/$(BIN) Dockerfile.in
	@sed \
	    -e 's|ARG_BIN|$(BIN)|g' \
	    -e 's|ARG_ARCH|$(ARCH)|g' \
	    -e 's|ARG_FROM|$(BASEIMAGE)|g' \
	    Dockerfile.in > .dockerfile-$(ARCH)
	@docker build -t $(IMAGE):$(VERSION) -f .dockerfile-$(ARCH) .
	@docker images -q $(IMAGE):$(VERSION) > $@

container-name:
	@echo "container: $(IMAGE):$(VERSION)"

push: .push-$(DOTFILE_IMAGE) push-name
.push-$(DOTFILE_IMAGE): .container-$(DOTFILE_IMAGE)
ifeq ($(findstring gcr.io,$(REGISTRY)),gcr.io)
	@gcloud docker -- push $(IMAGE):$(VERSION)
else
	@docker push $(IMAGE):$(VERSION)
endif
	@docker images -q $(IMAGE):$(VERSION) > $@

push-name:
	@echo "pushed: $(IMAGE):$(VERSION)"

version:
	@echo $(VERSION)

test: build-dirs
	@docker run                                                            \
	    -ti                                                                \
	    -u $$(id -u):$$(id -g)                                             \
	    -v $$(pwd)/.go:/go:$(DOCKER_MOUNT_MODE)                            \
	    -v $$(pwd)/$(CMD):/go/src/$(PRO)/$(CMD):$(DOCKER_MOUNT_MODE)                     \
        -v $$(pwd)/$(PKG):/go/src/$(PRO)/$(PKG):$(DOCKER_MOUNT_MODE)                     \
        -v $$(pwd)/$(DEP):/go/src/$(PRO)/$(DEP):$(DOCKER_MOUNT_MODE)                     \
        -v $$(pwd)/bin/$(ARCH):/go/bin:$(DOCKER_MOUNT_MODE)                \
        -v $$(pwd)/bin/$(ARCH):/go/bin/linux_$(ARCH):$(DOCKER_MOUNT_MODE)  \
        -v $$(pwd)/.go/std/$(ARCH):/usr/local/go/pkg/linux_$(ARCH)_static:$(DOCKER_MOUNT_MODE)  \
        -w /go/src                                                 \
        $(BUILD_IMAGE)                                                     \
        /bin/sh -c "                                                       \
	        ./build/test.sh                                                \
	    "

build-dirs:
	@mkdir -p bin/$(ARCH)
	@mkdir -p .go/src/$(PKG) .go/pkg .go/bin .go/std/$(ARCH)

clean: container-clean bin-clean

container-clean:
	rm -rf .container-* .dockerfile-* .push-*

bin-clean:
	rm -rf .go bin
