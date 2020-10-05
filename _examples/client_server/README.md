### Client-Server Example

To generate client-side code using dbus-codegen-go, run
 ```dbus-codegen-go -package=main -output=client/client.go -client-only example.xml```
 
To generate server-side code using dbus-codegen-go, run
 `dbus-codegen-go -package=main -output=server/server.go -server-only example.xml`
 
Now, *client.go and server.go files for the contract defined in example.xml will be generated in client/ and server/ folders respectively.*
 
To build go binaries for both client and server, run the below command inside client/ and server/ respectively
 `go build .`
 
In parallel consoles, run 
``` 
./server
Listening on interface - org.example.Demo and path /org/example/Demo ...
```

``` 
./client
Message over dBus : Hello, this is example codebase for dbus-codegen-go
```
#### Note
Make sure dbus-codegen-go is installed. To install, run
`GO111MODULE=on go get -u github.com/amenzhinsky/dbus-codegen-go`

