VERSION = v0.1.0

test:
	go test -v ./...
build:
	go build -o out/check ./cmd/check && \
 	go build -o out/in ./cmd/in
build-linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o out/check ./cmd/check && \
 	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o out/in ./cmd/in
docker:
	docker build -t mtharrison/github-pr-comment-resource:$(VERSION) . && \
	docker push mtharrison/github-pr-comment-resource:$(VERSION)
