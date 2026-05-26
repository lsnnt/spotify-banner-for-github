# spotify banner

It's an http server running on port 8080 which get's your 20 recently played songs and then converts it into svg then returns it.


## To have your own banner

First of all clone and go to the directory

1) Create a .env file
get the sp_dc cookie value from the cookies section on open.spotify.com
```
SPDC="Your sp_dc cookie value"
```
2) Build the binary file
```
go build . 
```
3) Run the binary file
```
./spotify-banner-for-github
```

Then go to http://localhost:8080/ to view the svg image generated


This is the initial design, clone the repo contribute it, improve it

we do need improvements on svg design

## Banner 
<img src="https://spotify.fustin.top/?temp=foo" />
