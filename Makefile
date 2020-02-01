# Run tests
test: fmt vet
	go test ./pkg/... ./cmd/... -coverprofile cover.out

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# build Binary
build:
	go build ./cmd/kubectl-pexec.go

# all
all: test build