## Logging usage

```go
package main

import (
        "fmt"
        "time"
        "github.com/philips-software/go-hsdp-api/logging"
)

func main() {
        client, err := logging.NewClient(http.DefaultClient, logging.Config{
                SharedKey:    "YourSharedKey",
                SharedSecret: "YourSharedSecret",
                BaseURL:      "https://logingestor-xxx.foo.com",
                ProductKey:   "YourProductKey",
        })
        var logResource = logging.Resource{
           ID: "856b1142-6df5-4c84-b11d-da3f0a794e84",
           EventID: "1",
           Category: "ApplicationLog",
           Component: "TestApp",
           TransactionID: "1f12f95c-77a0-48da-835d-e95aa116198f", // traceability
           ServiceName: "TestApp",
           ApplicationInstance: "7248e79e-ba0b-4d0e-82a9-fb7a47d26c23",
           OriginatingUser: "729e83bb-ce7d-4052-92f8-077a376d774c",
           Severity: "Info",
           LogTime: time.Now().Format("2006-01-02T15:04:05.000Z07:00")
           LogData.Message: "Test log message",
        }
        _, err := client.StoreResource([]logging.Resource{ logResource }, 1)
        if err != nil {
            fmt.Printf("Batch flushing failed: %v\n", err)
        }
}
```

## TODO

- Increase API coverage
- Increase code coverage

## Issues

- If you have an issue: report it on the [issue tracker](https://github.com/hsdp/go-hsdp-api/issues)

## Author

Andy Lo-A-Foe (<andy.lo-a-foe@philips.com>)

## License

License is MIT. See [LICENSE file](LICENSE.md)
