
.PHONY: serve


serve:
	env GOOGLE_APPLICATION_CREDENTIALS=credentials.json SECRET_KEY=2kkcidiwkkkwm PASSWORD=1 go run main.go -dev