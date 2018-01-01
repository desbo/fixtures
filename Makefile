build:
	swagger generate server --exclude-main --flag-strategy=pflag
	spectacle -t docs swagger.yaml

deploy:
	gcloud app deploy ./app/app.yaml --project tt365-fixtures
