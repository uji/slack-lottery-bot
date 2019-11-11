.PHONY: deps clean build

deps:
	go get -u ./...

clean:
	rm -rf ./select/select
	rm -rf ./lottery/lottery

build:
	GOOS=linux GOARCH=amd64 go build -o select/select ./select
	GOOS=linux GOARCH=amd64 go build -o lottery/lottery ./lottery

export:
	export AWS_PROFILE=default

package:
	sam package --template-file template.yaml --output-template-file output-template.yaml --s3-bucket slack-lottery-bot

deploy:
	sam deploy --template-file output-template.yaml --stack-name slack-lottery-bot --capabilities CAPABILITY_IAM
