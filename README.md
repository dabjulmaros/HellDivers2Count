# HEllDIVERS 2 Steam Count

[View the site](https://helldivers2count-production.up.railway.app/)

This is a crude approximation of what SteamDB offers focused on the HELLDIVERS 2 Game.

This was initially created while staring and the `Servers at capacity` error, and later refined, cause I'm really enjoying working with GO.

This project uses HTMX for request. GO as the server, and some basic templates.

Some JS was unavoidable for the charts, and for transforming the HTMX request into data that the charts library could use.

## Uses

[Steam CurrentPlayers API](https://partner.steamgames.com/doc/webapi/ISteamUserStats#GetNumberOfCurrentPlayers)

[Charts.js](https://www.chartjs.org/)

[HTMX](https://htmx.org/)

[GO](https://go.dev/)

## Deploying Locally

If you already have go installed, cloning this project and running `go run main.go` should be enough.
