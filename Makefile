clean:
	rm -rf ./dist
	docker rmi slack-lottery-bot_cdk

build:
	mkdir -p ./dist/select
	mkdir -p ./dist/lottery
	GOOS=linux GOARCH=amd64 go build -o dist/select/bin ./select
	GOOS=linux GOARCH=amd64 go build -o dist/lottery/bin ./lottery
	zip -r ./dist/select.zip ./dist/select
	zip -r ./dist/lottery.zip ./dist/lottery
	rm -rf ./dist/select ./dist/lottery

deploy-by-docker:
	docker-compose build --no-cache
	make build
	make up
	docker-compose exec cdk npm run build
	docker-compose exec cdk cdk synth
	docker-compose exec cdk cdk bootstrap
	docker-compose exec cdk cdk deploy --require-approval never
	make down

up:
	docker-compose up -d

down:
	docker-compose down
