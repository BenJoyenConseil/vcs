
.PHONY: test
test:
	go test -v ./...

.PHONY: install
install:
	go install 

.PHONY: clean
clean:
	rm -Rf *.ugit vcs

.PHONY: build
build:
	go build -v .
	@echo "[INFO] Binaries availlable : "
	@ls -l vcs