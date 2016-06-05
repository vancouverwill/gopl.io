# Go FTP Server and Client

For the example it is better to have the server and client running in seperate directories so we can see that the client is workig and updating or returning information about the server directory

```
cd ./ftp_server && go build -o ftp_server  ./main.go && ./ftp_server 8000 & 
cd ./ftp_client && go build -o ftp_client ./main.go && ./ftp_client 8000
```
