# notifications design

Server side notification mechanism.

## workflow

1.  websocket client connects and gets stored to list
  1.  if websocket client disconnects, remove from list
1.  client can control which page gets shown by issuing command to server
  1.  notification gets sent to all pages
  1.  routing on client side redirects page
  1.  clients should be listening at all times for every page

## manager protocol

`/ws/manager`

manager can control what gets displayed at what time

### WS_LOGIN

TODO: credentials should be passed in such that malicious users cannot take control

### SCOREBOARD

### ADVERTISEMENT

### VIDEO

### PHOTO


## subscriber protocol

`/ws/subscriber`

### TITLE/SETUP

### SCOREBOARD

### ADVERTISEMENT

### VIDEO

### PHOTO