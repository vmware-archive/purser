# Copyright (c) 2018 VMware Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#    http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The binary to build (just the basename).
BIN := controller

# This repo's root import path (under GOPATH)
PRO := github.com/vmware/purser
DEP := vendor
BUILD := build
PKG := pkg
CMD := cmd/controller

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

# Set dep management tool parameters
VENDOR_DIR := vendor
DEP_BIN_NAME := dep
DEP_BIN_DIR := ./tmp/bin
DEP_BIN := $(DEP_BIN_DIR)/$(DEP_BIN_NAME)
DEP_VERSION := v0.5.0

# Define and get the vakue for UNAME_S variable from shell
UNAME_S := $(shell uname -s)


.PHONY: travis-build
travis-build: install-plugin install-controller travis-success

.PHONY: install-plugin
install-plugin:
	go install github.com/vmware/purser/cmd/plugin

.PHONY: install-controller
install-controller: build container

.PHONY: travis-success
travis-success:
	@echo "travis build success"

# If you want to build all binaries, see the 'all-build' rule.
# If you want to build all containers, see the 'all-container' rule.
# If you want to build AND push all containers, see the 'all-push' rule.
.PHONY: all
all: deps build check

build-%:
	@$(MAKE) --no-print-directory ARCH=$* build

container-%:
	@$(MAKE) --no-print-directory ARCH=$* container

push-%:
	@$(MAKE) --no-print-directory ARCH=$* push

.PHONY: all-build
all-build: $(addprefix build-, $(ALL_ARCH))

.PHONY: all-container
all-container: $(addprefix container-, $(ALL_ARCH))

.PHONY: all-push
all-push: $(addprefix push-, $(ALL_ARCH))

.PHONY: build
build: bin/$(ARCH)/$(BIN)

bin/$(ARCH)/$(BIN): build-dirs
	@echo "building: $@"
	@docker run                                                            \
	    -ti                                                                \
	    -u $$(id -u):$$(id -g)                                             \
	    -v $$(pwd)/.go:/go:$(DOCKER_MOUNT_MODE)                            \
        -v $$(pwd)/$(BUILD):/go/src/$(PRO)/$(BUILD):$(DOCKER_MOUNT_MODE)   \
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
	        ./$(PRO)/$(BUILD)/build.sh                                               \
	    "

DOTFILE_IMAGE = $(subst :,_,$(subst /,_,$(IMAGE))-$(VERSION))

.PHONY: container
container: .container-$(DOTFILE_IMAGE) container-name
.container-$(DOTFILE_IMAGE): bin/$(ARCH)/$(BIN) Dockerfile.in
	@sed \
	    -e 's|ARG_BIN|$(BIN)|g' \
	    -e 's|ARG_ARCH|$(ARCH)|g' \
	    -e 's|ARG_FROM|$(BASEIMAGE)|g' \
	    Dockerfile.in > .dockerfile-$(ARCH)
	@docker build -t $(IMAGE):$(VERSION) -f .dockerfile-$(ARCH) .
	@docker images -q $(IMAGE):$(VERSION) > $@

.PHONY: container-name
container-name:
	@echo "container: $(IMAGE):$(VERSION)"

.PHONY: push
push: .push-$(DOTFILE_IMAGE) push-name
.push-$(DOTFILE_IMAGE): .container-$(DOTFILE_IMAGE)
ifeq ($(findstring gcr.io,$(REGISTRY)),gcr.io)
	@gcloud docker -- push $(IMAGE):$(VERSION)
else
	@docker push $(IMAGE):$(VERSION)
endif
	@docker images -q $(IMAGE):$(VERSION) > $@

.PHONY: push-name
push-name:
	@echo "pushed: $(IMAGE):$(VERSION)"

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: test
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

.PHONY: build-dirs
build-dirs:
	@mkdir -p bin/$(ARCH)
	@mkdir -p .go/src/$(PKG) .go/pkg .go/bin .go/std/$(ARCH)

.PHONY: clean
clean: container-clean bin-clean

.PHONY: container-clean
container-clean:
	rm -rf .container-* .dockerfile-* .push-*

.PHONY: bin-clean
bin-clean:
	rm -rf .go bin

.PHONY: clean-vendor
## Removes the ./vendor directory.
clean-vendor:
	-rm -rf $(VENDOR_DIR)

.PHONY: deps
## Download build dependencies.
deps: $(DEP_BIN) $(VENDOR_DIR)

# install dep in a the tmp/bin dir of the repo
$(DEP_BIN):
	@echo "Installing 'dep' $(DEP_VERSION) at '$(DEP_BIN_DIR)'..."
	mkdir -p $(DEP_BIN_DIR)
ifeq ($(UNAME_S),Darwin)
	@curl -L -s https://github.com/golang/dep/releases/download/$(DEP_VERSION)/dep-darwin-amd64 -o $(DEP_BIN)
	@cd $(DEP_BIN_DIR) && \
	echo "1a7bdb0d6c31ecba8b3fd213a1170adf707657123e89dff234871af9e0498be2  dep" > dep-darwin-amd64.sha256 && \
	shasum -a 256 --check dep-darwin-amd64.sha256
else
	@curl -L -s https://github.com/golang/dep/releases/download/$(DEP_VERSION)/dep-linux-amd64 -o $(DEP_BIN)
	@cd $(DEP_BIN_DIR) && \
	echo "287b08291e14f1fae8ba44374b26a2b12eb941af3497ed0ca649253e21ba2f83  dep" > dep-linux-amd64.sha256 && \
	sha256sum -c dep-linux-amd64.sha256
endif
	@chmod +x $(DEP_BIN)

$(VENDOR_DIR): Gopkg.toml Gopkg.lock
	@echo "checking dependencies..."
	@$(DEP_BIN) ensure -v

GOFORMAT_FILES := $(shell find  . -name '*.go')

.PHONY: format ## Formats any go file that differs from gofmt's style and removes unused imports
format: 
	@gofmt -s -l -w ${GOFORMAT_FILES}
	@goimports -l -w ${GOFORMAT_FILES}

.PHONY: tools
tools: ## Installs required go tools
	@go get -u github.com/alecthomas/gometalinter && gometalinter --install
	@go get -u golang.org/x/tools/cmd/goimports
	
.PHONY: install
install: ## Fetches all dependencies using dep
	@$(DEP_BIN) ensure -v

.PHONY: update
update: ## Updates all dependencies defined for dep
	@$(DEP_BIN) ensure -update -v

.PHONY: check
check: ## Concurrently runs a whole bunch of static analysis tools
	gometalinter --enable=misspell --enable=gosimple --enable-gc --vendor --deadline 300s ./...	

include ./.make/Makefile.deploy	