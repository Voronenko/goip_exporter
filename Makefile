VERSION=`cat version.txt`
SHORTSHA=`git rev-parse --short HEAD`

LDFLAGS=-X main.appVersion=$(VERSION)
LDFLAGS+=-X main.shortSha=$(SHORTSHA)

build:
	mkdir -p dist
	go build -o dist/goip_exporter -ldflags "$(LDFLAGS)" .

build-arm:
	mkdir -p dist
	GOOS=linux GOARCH=arm go build -o dist/goip_exporter -ldflags "$(LDFLAGS)" .

utils:
	go get github.com/mitchellh/gox
	go get github.com/tcnksm/ghr

update-vendor:
#	rm -rf vendor
	go mod vendor

snapshot:
	goreleaser release --rm-dist --snapshot

prerelease:
	goreleaser release --rm-dist --skip-publish --skip-sign

docker-arm-build:
	docker build -f Dockerfile.arm -t voronenko/goip_exporter:$(VERSION)-pi .

docker-arm-push:
	docker build -f Dockerfile.arm -t voronenko/goip_exporter:$(VERSION)-pi .
