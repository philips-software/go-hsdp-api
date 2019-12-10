## Using the TPNS client

```go
package main

import (
        "fmt"
        "time"
        "github.com/philips-software/go-hsdp-api/tpns"
)

func main() {
        client, err := tpns.NewClient(&tpns.Config{
                Username: "Foo",
                Password: "YourP@ssword!",
                TPNSURL:  "https://tpns.foo.com",
        })
        if err != nil {
            fmt.Printf("Error creating client: %v\n", err)
            return
        } 
        ok, resp, err := tpns.Messages.Push(&tpns.Message{
            Content:       "YAY! It is working!",
            PropositionID: "XYZ",
            MessageType:   "Push",
            Targets:       []string{"5b78e5f8-d73f-4712-aae0-355f6fa91752"},
        })
        if err != nil {
            fmt.Printf("Error pushing: %v\n", err)
            return
        }
        if !ok {
            fmt.Printf("Error pushing: %d\n", resp.StatusCode)
            return
        }
        fmt.Printf("Push success: %d\n", resp.StatusCode)
}
```
