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


## notes
clissold club ID: 5123

central london division IDs:
	1: 5596,
	2: 5597,
	3: 5598,
	4: 5599,
	5: 5600,
	6: 5601

## example URLs:
all clissold TTC games: http://localhost:8080/CentralLondon/Winter_2017-18?club_id=5123
clissold TTC division 6: http://localhost:8080/CentralLondon/Winter_2017-18?club_id=5123&division_id=5601
