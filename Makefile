# =================================================================
#
# Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================
bin/gocat_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/gocat_darwin_amd64 $$(go list ./...)

bin/gocat_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/gocat_linux_amd64 $$(go list ./...)

bin/gocat_windows_amd64.exe:
	GOOS=windows GOARCH=amd64 go build -o bin/gocat_windows_amd64.exe $$(go list ./...)

bin/gocat_linux_arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/gocat_linux_arm64 $$(go list ./...)

build: \
bin/gocat_darwin_amd64 \
bin/gocat_linux_amd64 \
bin/gocat_windows_amd64.exe \
bin/gocat_linux_arm64

fmt:
	go fmt $$(go list ./... )

install:
	go install $$(go list ./...)

test:
	bash scripts/test.sh

clean:
	rm -fr bin
