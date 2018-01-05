.PHONY: check test lint vet fmt-check

PKGS := $(shell go list ./... | grep -v /vendor/ | grep -v loraserver/api | grep -v /migrations | grep -v /static)

check: lint vet fmt-check test

test:
	@go test -p 1 -v $(PKGS)

lint:
	@for pkg in $(PKGS) ; do \
		! echo $$pkg | grep '/internal' > /dev/null; \
		SKIP=$$? ; \
		if [ $$SKIP -eq 1 ]; then continue; fi ; \
		golint -set_exit_status $$pkg ; \
		EXIT_CODE=$$? ; \
		if [ $$EXIT_CODE -ne 0 ]; then exit 1; fi ; \
	done

vet:
	@go vet $(PKGS)

fmt-check:
	@ gofmt -l . | grep -v ^vendor/; \
	EXIT_CODE=$$?; \
	if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi

