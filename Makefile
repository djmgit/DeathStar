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
prepare_for_linux=(export GOOS=linux)
prefix=/usr

.PHONY: all build clean install uninstall build_linux prepare_linux

all: build

build:	clean
		@echo "building ${BINARY_NAME}..."
		@mkdir $(DIST_DIR)
		$(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME) .
		@echo "build successful"
		@echo "DeathStar binary can be found inside ${DIST_DIR}"

clean:
		$(GOCLEAN)
		@rm -rf $(DIST_DIR)
		@echo "build cleaned"

lambda_clean:
		@rm $(LAMBDA_PACKAGE)

install:
		@install -D $(DIST_DIR)/$(BINARY_NAME) $(DESTDIR)$(prefix)/bin/$(BINARY_NAME)
		@echo "installed DeathStar"

lambda_package: lambda_clean
		@echo "Create lambda zip package"
		cd $(DIST_DIR); zip $(LAMBDA_PACKAGE) $(BINARY_NAME)
		@mv $(DIST_DIR)/$(LAMBDA_PACKAGE) ./

uninstall:
		@rm -f $(DESTDIR)$(prefix)/bin/$(BINARY_NAME)
		@echo "uinsalled DeathStar"

doc:
		@$(GODOC) -http=:6060
