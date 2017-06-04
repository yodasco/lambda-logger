deploy:
	glide install
	apex deploy \
	-s GOOGLE_APPLICATION_CREDENTIALS=credentials.json \
	-s GOOGLE_PROJECT_ID=`jq -r .project_id functions/stackdriver/credentials.json`

logs:
	apex logs -f

invoke-remote:
	apex invoke lambda-logger < functions/stackdriver/event.json

run:
	glide install
	GOOGLE_APPLICATION_CREDENTIALS=functions/stackdriver/credentials.json \
	go run main.go


setup:
	@echo Install Go
	@echo Install Apex
	@echo Install Glide
	@echo Install jq
