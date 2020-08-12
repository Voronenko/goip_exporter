update-vendor:
#	rm -rf vendor
	go mod vendor
build:
	mkdir -p dist
	go build -o dist/goip-exporter main.go
snapshot:
	goreleaser release --rm-dist --snapshot
prerelease:
	goreleaser release --rm-dist --skip-publish --skip-sign
