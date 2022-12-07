# Go(rilla) websocket server example with SQL DB for persistent user-data and sample client.

## serves http by default on port :8080

number of maximum connections depends on OS. max for windows is 2000. linux can be increased further. windows default is between 256 and 1028. 

### connection to valid (my)SQL DB is required. (dialect is exchangeable)

after a connection is promoted to websocket, the client is asked for credentials. -> "request_credentials"

client responds with credentials. <- "credentials"

-in case of valid credentials, server responds with saved user data (as JSON). -> "user_data"

-in case of invalid credentials, server responds with new user data (as JSON). -> "user_data"

### the websocket connection is now linked to a persistent unique user. the server keeps track of all sockets that belong to each user

highly expandable
