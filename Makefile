# ==================== [START] Global Variable Declaration =================== #
SHELL := /bin/bash
BASE_DIR := $(shell pwd)
UNAME_S := $(shell uname -s)
APP_NAME := decker

export
# ===================== [END] Global Variable Declaration ==================== #

# =========================== [START] git scripts ============================ #
init:
	@git config core.hooksPath githooks
# ============================ [END] git scripts ============================= #

# =========================== [START] Build Scripts ========================== #
clean:
	@find $(BASE_DIR)/internal/app/decker/plugins -name "*.so" -type f -delete
	@rm -rf $(BASE_DIR)/vendor $(BASE_DIR)/_vendor-*
	@rm -rf ./stash.sqlite
	@rm -rf ./wafw00f.egg-info

build_plugins:
	@$(BASE_DIR)/scripts/build-plugins.sh $(BASE_DIR)/internal/app/decker/plugins

build_plugin:
	@$(BASE_DIR)/scripts/build-plugins.sh $(BASE_DIR)/internal/app/decker/plugins $(plugin)

build_decker:
	@cd $(BASE_DIR)/cmd/decker && \
		echo "Building" $(APP_NAME) && \
		CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o $(BASE_DIR)/$(APP_NAME)

build_all: build_plugins build_decker

docker_build:
	@docker build -f ./build/package/kali.Dockerfile -t stevenaldinger/$(APP_NAME):latest .

docker_build_prod:
	@docker build -f ./build/package/kali.Dockerfile -t stevenaldinger/$(APP_NAME):latest .

docker_build_kali:
	@docker build -f ./build/package/kali.Dockerfile -t stevenaldinger/$(APP_NAME):kali .

docker_build_minimal:
	@docker build -f ./build/package/minimal.Dockerfile -t stevenaldinger/$(APP_NAME):minimal .

docker_build_and_push: docker_build_minimal docker_build_prod docker_build_kali
	@docker tag stevenaldinger/$(APP_NAME):minimal stevenaldinger/$(APP_NAME):minimal-$(tag)
	@docker tag stevenaldinger/$(APP_NAME):latest stevenaldinger/$(APP_NAME):base-$(tag)
	@docker tag stevenaldinger/$(APP_NAME):kali stevenaldinger/$(APP_NAME):kali-$(tag)
	@docker push stevenaldinger/$(APP_NAME):minimal
	@docker push stevenaldinger/$(APP_NAME):minimal-$(tag)
	@docker push stevenaldinger/$(APP_NAME):latest
	@docker push stevenaldinger/$(APP_NAME):base-$(tag)
	@docker push stevenaldinger/$(APP_NAME):kali
	@docker push stevenaldinger/$(APP_NAME):kali-$(tag)

# ============================ [END] Build Scripts =========================== #

# ============================ [START] Run Scripts =========================== #
run:
	@$(BASE_DIR)/$(APP_NAME) ./examples/example.hcl

run_hello_world:
	@$(BASE_DIR)/$(APP_NAME) ./examples/hello-world.hcl

docker_run:
	@echo "Forwarding port 6060 for godoc usage within the container."
	@docker run -it --rm \
		-v $(BASE_DIR):/go/src/github.com/stevenaldinger/$(APP_NAME) \
		-v $(HOME)/decker-reports:/tmp/reports \
		-p 6060:6060 \
	 stevenaldinger/$(APP_NAME):kali bash

docker_run_prod:
	@docker run -it --rm \
		stevenaldinger/$(APP_NAME):latest

docker_run_dvwa:
	@docker-compose -f $(BASE_DIR)/deployments/docker-compose-dvwa.yml up -d
# ============================= [END] Run Scripts ============================ #

# =========================== [START] Stop Scripts =========================== #
docker_stop_dvwa:
	@docker-compose -f $(BASE_DIR)/deployments/docker-compose-dvwa.yml down
# ============================ [END] Stop Scripts ============================ #

# ========================= [START] Formatting Script ======================== #
gofmt:
	@go fmt github.com/stevenaldinger/decker/...

golint:
	@golint github.com/stevenaldinger/decker/cmd/...
	@golint github.com/stevenaldinger/decker/internal/...

govet:
	@go vet github.com/stevenaldinger/decker/cmd/...
	@go vet github.com/stevenaldinger/decker/internal/...

lint: gofmt golint govet
# ========================== [END] Formatting Script ========================= #
test_by_pkg:
	@go test -v github.com/stevenaldinger/decker/internal/pkg/dependencies
	@go test -v github.com/stevenaldinger/decker/internal/pkg/hcl
	@go test -v github.com/stevenaldinger/decker/internal/pkg/paths
	@go test -v github.com/stevenaldinger/decker/internal/pkg/plugins
	@go test -v github.com/stevenaldinger/decker/internal/pkg/reports

cvg_by_pkg:
	@go test -v -cover github.com/stevenaldinger/decker/internal/pkg/dependencies
	@go test -v -cover github.com/stevenaldinger/decker/internal/pkg/hcl
	@go test -v -cover github.com/stevenaldinger/decker/internal/pkg/paths
	@go test -v -cover github.com/stevenaldinger/decker/internal/pkg/plugins
	@go test -v -cover github.com/stevenaldinger/decker/internal/pkg/reports

gotest: test_by_pkg
	@cd $(BASE_DIR)/cmd/decker && \
		go test -v -cover
# ======================= [START] Documentation Scripts ====================== #
godoc:
	@godoc -http=":6060"
# ==============-========= [END] Documentation Scripts =========-============= #
