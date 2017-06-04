deploy:
	glide install
	apex deploy \
	-s GOOGLE_APPLICATION_CREDENTIALS=credentials.json \
	-s GOOGLE_PROJECT_ID=yodas-test

logs:
	apex logs -f

invoke-remote:
	apex invoke lambda-logger < functions/lambda-logger/event.json

run:
	glide install
	GOOGLE_APPLICATION_CREDENTIALS=functions/lambda-logger/credentials.json \
	go run main.go
