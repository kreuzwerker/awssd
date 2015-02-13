TOKEN = `cat .token`
REPO := awssd
USER := kreuzwerker
VERSION := "v0.0.1"

build:
	mkdir -p out/darwin out/linux
	GOOS=darwin go build -o out/darwin/awssd -ldflags "-X main.build `git rev-parse --short HEAD`" bin/awssd.go
	GOOS=linux go build -o out/linux/awssd -ldflags "-X main.build `git rev-parse --short HEAD`" bin/awssd.go

clean:
	rm -rf out

release: clean build
	github-release release --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN)
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name $(REPO)-osx --file out/darwin/$(REPO)
	github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name $(REPO)-linux --file out/linux/$(REPO)

test:
	go test -cover

install:
	go get github.com/awslabs/aws-sdk-go/aws
	go get github.com/awslabs/aws-sdk-go/gen/ec2
	go get github.com/awslabs/aws-sdk-go/gen/route53
	go get github.com/deckarep/golang-set
	go get github.com/clipperhouse/gen
	go get github.com/clipperhouse/set
	go get github.com/stretchr/testify/assert
	gen
