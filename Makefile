# include the common makefile
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))


ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell pwd -P))
endif

# Linux command settings
CODE_DIRS := $(ROOT_DIR)/internal $(ROOT_DIR)/cmd $(ROOT_DIR)/tests $(ROOT_DIR)/plugin $(ROOT_DIR)/plugins
FIND := find $(CODE_DIRS)

format:
	@echo "===========> Formating codes"
	@$(FIND) -type f -name '*.go' | xargs gofmt -s -w
	@$(FIND) -type f -name '*.go' | xargs goimports -l -w -local $(ROOT_PACKAGE)
