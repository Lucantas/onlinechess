- Create the hub
    - Create the socket based on the responseWriter, request and hub
- Associate the local address
- Wait for incoming connections from clients
- Accept incoming connection
    - Add every user to a lobby
    - Generate uuid for every user that connected to the lobby
    - randomly add two users identified by a uuid to a room and removing them from the waiting room
- communicate with client
- close the soccket descriptor
    - If a user disconnect from a match
        - the other user receives the signal and might choose to start another game