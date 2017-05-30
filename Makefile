deploy:
	apex deploy \
		-s GOOGLE_APPLICATION_CREDENTIALS=credentials.json

logs:
	apex logs -f

invoke-remote:
	apex invoke lambda-logger < functions/lambda-logger/event.json

run:
	GOOGLE_APPLICATION_CREDENTIALS=functions/lambda-logger/credentials.json \
	go run main.go
