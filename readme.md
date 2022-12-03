minimalistic Go(rilla) websocket server with minimalistic react frontend sample. 

after a connection is promoted to websocket, the client is asked for credentials. -> "request_credentials"

client responds with credentials. <- "credentials"

-in case of valid credentials, server responds with saved user data. -> "user_data"

-in case of invalid credentials, server responds with new user data. -> "user_data"

highly expandable
