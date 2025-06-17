#
# Display help.
#
.PHONY: help
help:
	@echo 'Available tasks:'
	@echo '  make help     -- Display this help message'
	@echo '  make build    -- Build the `check_cloudwatch` binary'
	@echo '  make clean    -- Remove build artifacts'
	@echo '  make test     -- Run tests'
	@echo '  make license  -- Collect license information and save it to `./licenses/`'
	@echo '  make release  -- Release the check_cloudwatch binary'

#
# Build the binary.
#
.PHONY: build
build: clean
	@echo 'Building check_cloudwatch binary...'

	go build \
		-v \
		./cmd/check_cloudwatch/

#
# Remove build artifacts.
#
.PHONY: clean
clean:
	@echo 'Cleaning build artifacts...'

	rm -rf \
	  ./check_cloudwatch \
	  ./dist/ \
	  ./licenses/

#
# Run tests.
#
.PHONY: test
test:
	@echo 'Running tests...'

	go test -v ./...

#
# Collect license information.
#
.PHONY: license
license:
	@echo 'Collecting license information...'

	go tool go-licenses save ./cmd/check_cloudwatch/ \
	  --force \
	  --save_path ./licenses/

#
# Release the binary.
#
.PHONY: release
release: clean license
	@echo 'Releasing check_cloudwatch binary...'

	go tool goreleaser release --clean
