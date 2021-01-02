GOCMD=go
GODOC=godoc
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=deathstar
LAMBDA_PACKAGE=deathstar.zip
DIST_DIR=dist
CONF_DIR=conf
prefix=/usr

.PHONY: all build clean install uninstall build_linux prepare_linux

all: build

build:	clean
		@echo "building ${BINARY_NAME}..."
		@mkdir $(DIST_DIR)
		$(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME) .
		@echo "build successful"
		@echo "DeathStar binary can be found inside ${DIST_DIR}"

prepare_linux:
		@echo "preparing for linux build..."
		@export GOOS=linux

build_linux: prepare_linux build

clean:
		$(GOCLEAN)
		@rm -rf $(DIST_DIR)
		@echo "build cleaned"

install:
		@install -D $(DIST_DIR)/$(BINARY_NAME) $(DESTDIR)$(prefix)/bin/$(BINARY_NAME)
		@echo "installed DeathStar"

lambda_package:
		@echo "Create lambda zip package"
		@zip $(LAMBDA_PACKAGE) $(DIST_DIR)/$(BINARY_NAME)

uninstall:
		@rm -f $(DESTDIR)$(prefix)/bin/$(BINARY_NAME)
		@echo "uinsalled DeathStar"

doc:
		@$(GODOC) -http=:6060
