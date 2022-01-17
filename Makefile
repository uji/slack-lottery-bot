clean:
	make down
	docker rmi slack-lottery-bot_cdk

build:
	rm -rf ./dist
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/select ./select
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/lottery ./lottery

docker-cp:
	docker cp . slack-lottery-bot-cdk:/workdir

deploy-by-docker:
	make build
	make docker-cp
	docker-compose exec cdk npm run build
	docker-compose exec cdk cdk synth
	docker-compose exec cdk cdk bootstrap
	docker-compose exec cdk cdk deploy --require-approval never

up:
	docker-compose up -d

down:
	docker-compose down
