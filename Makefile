deploy:	glide-install
	apex deploy \
	-s GOOGLE_APPLICATION_CREDENTIALS=credentials.json

logs:
	apex logs -f

invoke-remote:
	apex invoke lambda-logger < functions/stackdriver/event.json

setup:
	@echo Install Go
	@echo Install Apex
	@echo Install Glide

glide-install:
	glide install
