# tabeltennis365.com API

## setup
1. install https://github.com/go-swagger/go-swagger
2. run `swagger generate server --exclude-main --flag-strategy=pflag` in this directory

`--exclude-main` is provided as AppEngine will handle running the app.
`--flag-strategy=pflag` is used as the alternative (`go-flags`) doesn't work with AppEngine.

## example URL structure
https://www.tabletennis365.com/CentralLondon/Fixtures/Winter_2017-18/All_Divisions
?c=False // ???
&vm=1 // view mode
&d= // division
&vn= // venue
&cl= // club
&t= // ???
&swn=True // show week numbers
&hc=False // hide completed