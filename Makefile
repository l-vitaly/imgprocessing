#############################
INSTALL_COMMON_VERSION=v1.0.2
#############################

CMD_PATH=./cmd/imgprocessing

-include common.mk

ifndef COMMON_VERSION
 $(shell git archive --remote=git@gitlab.com:skugrid-go/makefile-common.git $(INSTALL_COMMON_VERSION) | tar xvf - common.mk) \
 $(error "common.mk was fetched. please retry now")
endif

ifneq ($(COMMON_VERSION),$(INSTALL_COMMON_VERSION))
 $(shell git archive --remote=git@gitlab.com:skugrid-go/makefile-common.git $(INSTALL_COMMON_VERSION) | tar xvf - common.mk) \
 $(error "common.mk was updated. please retry now");
endif

wire:
	wire pkg/wire/wire.go