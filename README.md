# spotify banner

It fetches your 20 recently played songs and as of now prints it on terminal
The main reason to implement this by myself and not using and existing libraries is because spotify now restricts only premium spotify user to use the Web api so this is done using some tricks by logging it into your account and using developer try me web api thing.

## to use
First of all clone and go to the directory

1) Create a .env file
```
SPDC="Your sp_dc cookie value"
```
2) Run the script
```
go run .
```
To implement 

[x] An http server which then converts it into an svg to be included somewhere.
