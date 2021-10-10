build:
	swagger generate server --exclude-main --flag-strategy=pflag
	spectacle -t docs swagger.yaml

deploy:
	gcloud app deploy app.yaml --project planar-berm-308814

