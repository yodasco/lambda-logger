deploy:
	glide install
	apex deploy \
	-s GOOGLE_APPLICATION_CREDENTIALS=credentials.json \
	-s GOOGLE_PROJECT_ID=`jq -r .project_id functions/lambda-logger/credentials.json`

logs:
	apex logs -f

invoke-remote:
	apex invoke lambda-logger < functions/lambda-logger/event.json

run:
	glide install
	GOOGLE_APPLICATION_CREDENTIALS=functions/lambda-logger/credentials.json \
	go run main.go


setup:
	@echo Install Go
	@echo Install Apex
	@echo Install Glide
	@echo Install jq
