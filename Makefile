GOPATH=$(PWD)/.gopath
GO=GOPATH=$(GOPATH) go
GIT_SHA=$(shell git log -n 1 --pretty=format:%H)
VERSION=1.0
PROJECT_NAME=$(notdir $(basename $(PWD)))
DOCKER_USERNAME ?= docker_username

clean:
	rm -rf $(GOPATH)

init:
	#dependency management
	mkdir -p .gopath
	$(GO) get github.com/spf13/viper
	$(GO) get github.com/spf13/pflag
	$(GO) get github.com/op/go-logging
	$(GO) get github.com/go-sql-driver/mysql
	$(GO) get github.com/go-xorm/xorm
	$(GO) get github.com/stretchr/testify/assert
	$(GO) get github.com/pressly/chi
	$(GO) get github.com/auth0/go-jwt-middleware
	$(GO) get github.com/dgrijalva/jwt-go
	$(GO) get github.com/gophercloud/gophercloud
	$(GO) get github.com/gophercloud/gophercloud/openstack
	$(GO) get github.com/gophercloud/gophercloud/openstack/identity/v3/tokens
	$(GO) get github.com/mitchellh/mapstructure
	$(GO) get github.com/golang/mock/gomock
	$(GO) get github.com/golang/mock/mockgen
	$(GO) get github.com/kbhonagiri16/visualization-client
	$(GO) get github.com/xeipuuv/gojsonschema
	$(GO) get github.com/satori/go.uuid
	$(GO) get -u github.com/ulule/deepcopier
	$(GO) get -u gopkg.in/alecthomas/gometalinter.v1
	$(GO) get github.com/rubenv/sql-migrate/...
	GOPATH=$(GOPATH) $(GOPATH)/bin/gometalinter.v1 --install
	# as soon as our application does not use relative imports - source code
	# has to be present in GOPATH to make lint work
	# as soon as we created isolated GOPATH - we have to create a symlink
	# from GOPATH to our source code
	mkdir -p $(GOPATH)/src/$(PROJECT_NAME)/
	ln -s $(PWD) $(GOPATH)/src/$(PROJECT_NAME)

stylecheck:
	if [ ! -z "$$($(GO) fmt ./...)" ]; then exit 1; fi

fmt:
	$(GO) fmt ./...

lint:
	GOPATH=$(GOPATH) $(GOPATH)/bin/gometalinter.v1 --disable=gotype \
		  --disable=errcheck --disable=gas --disable=gocyclo --exclude=mock --exclude=Mock \
		  --exclude='dynamic type' ./...

generate-mocks:
	mkdir -p ./mock
	GOPATH=$(GOPATH) $(GOPATH)/bin/mockgen -destination ./mock/mock.go visualization-client ClientInterface,DatabaseManager,SessionInterface
	mkdir -p ./http_endpoint/common/mock
	GOPATH=$(GOPATH) $(GOPATH)/bin/mockgen -destination ./http_endpoint/common/mock/mock.go visualization-client/http_endpoint/common HandlerInterface,ClockInterface

clean-mocks:
	rm -r ./mock
	rm -r ./http_endpoint/common/mock

test: generate-mocks
	$(GO) test ./...

test-integration:
	docker run --name=grafana-integration-test -d -p 3000:3000 grafana/grafana
	sleep 10
	curl -v -X POST -H "Content-Type: applciation/json" -d '{"name":"PV Service", "login":"pv_service", "password":"123123"}' http://admin:admin@localhost:3000/api/admin/users
	curl -v -X PUT -H "Content-Type: applciation/json" -d '{"isGrafanaAdmin":true}' http://admin:admin@localhost:3000/api/admin/users/2/permissions
	GRAFANA_URL=http://0.0.0.0:3000 GRAFANA_USER=pv_service GRAFANA_PASS=123123 $(GO) test -v visualization-client/ -tags=integration	
	docker rm --force grafana-integration-test

build: fmt lint
	$(GO) build ./cmd/...

build-all: fmt
	mkdir -p build/linux-amd64
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-X main.version=$(VERSION) -X main.gitVersion=$(GIT_SHA)" -o $(PWD)/build/linux-amd64/visualizationapi ./cmd/visualizationapi
	GOOS=linux GOARCH=amd64 $(GO) build -o $(PWD)/build/linux-amd64/sql-migrate github.com/rubenv/sql-migrate/sql-migrate

package-init:
	docker build -t com.mirantis.pv/build ./tools/build

package-clean:
	docker image rm -f com.mirantis.pv/build
	rm -rf build/deb/*

package:
	docker run -e VERSION=$(VERSION) -v $(PWD):/app com.mirantis.pv/build /app/tools/build/build_deb.sh

package-debug:
	docker run -it -v $(PWD):/app com.mirantis.pv/build /bin/bash

docker:
	docker build -t $(DOCKER_USERNAME)/visualization-client -f tools/docker/visualization-client/Dockerfile .

docker-push:
	docker push $(DOCKER_USERNAME)/visualization-client

all: init fmt lint  build-all package-init package docker
