build := $(shell date +'%d/%m/%Y %H:%M')
flags := -s -w -X
# go build -race

build: export CGO_CFLAGS_ALLOW = -Xpreprocessor
build:
	$(eval version = $(shell bump2version --dry-run --allow-dirty --list patch | grep new_version | sed -r s,"^.*=",,))
	go build -ldflags="$(flags) 'main.version=$(version)' -X 'main.build=$(build)'"

patch:
	$(eval version = $(shell bump2version --allow-dirty --list patch | grep new_version | sed -r s,"^.*=",,))
	go build -ldflags="$(flags) 'main.version=$(version)' -X 'main.build=$(build)'"

minor:
	$(eval version = $(shell bump2version --allow-dirty --list minor | grep new_version | sed -r s,"^.*=",,))
	go build -ldflags="$(flags) 'main.version=$(version)' -X 'main.build=$(build)'"

major:
	$(eval version = $(shell bump2version --allow-dirty --list major | grep new_version | sed -r s,"^.*=",,))
	go build -ldflags="$(flags) 'main.version=$(version)' -X 'main.build=$(build)'"

win:
	$(eval version = $(shell bump2version --dry-run --allow-dirty --list patch | grep new_version | sed -r s,"^.*=",,))
	env GOOS=windows GOARCH=amd64 go build -ldflags="$(flags) 'main.version=$(version)' -X 'main.build=$(build)'"

clean:
	go clean
pack:
	upx --brute dispatcher