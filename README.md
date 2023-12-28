[![Go Report Card](https://goreportcard.com/badge/github.com/L1ghtn1ng/gowol)](https://goreportcard.com/report/github.com/L1ghtn1ng/gowol)

To install the Wake-on-LAN (WoL) module ("gowol"), you can use the go install command:
```go
go install github.com/L1ghtn1ng/gowol
```
This will compile and install the package and its dependencies to the $GOBIN directory.

To use the package in your Go code, you can import it like this:
```go
import "github.com/L1ghtn1ng/gowol"
```
Then, you can use the SendMagicPacket function to send a WoL magic packet to the target device like this:
```go
err := gowol.SendMagicPacket("00:11:22:33:44:55", "")
if err != nil {
	// Handle the error
}
```

You will need to replace the macAddress argument with the MAC address of the target device in the format ``"00:11:22:33:44:55"``. The ipAddress argument is optional, and if it is not specified or is an empty string, the default value of ``255.255.255.255`` will be used. You can specify a custom IP broadcast address by passing it as the ipAddress argument.

To run the tests for the ``"gowol"`` package, you can use the go test command:
```go
go test github.com/L1ghtn1ng/gowol
```
This will run all the tests in the "gowol" package and print the results.
