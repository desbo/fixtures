build:
	swagger generate server --exclude-main --flag-strategy=pflag
	spectacle -t docs swagger.yaml

