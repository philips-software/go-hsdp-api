# Using the Iron API client
HSDP uses Iron.io for container scheduling. This API client
implements a basic set of functionality to create code definitions and queue tasks.
The focus is on Docker code packages. 

# Registering a docker image
```go
package main

import (
        "fmt"
        "github.com/philips-software/go-hsdp-api/iron"
)

var (
    projectID = "yourIronProjectID"
    projectToken = "yourIronProjectToken"
    clusterID = "yourIronClusterID"
)

func main() {
        client, err := iron.NewClient(&iron.Config{
                ProjectID: projectID,
                Token:     projectToken,
                ClusterInfo: []iron.ClusterInfo{
                        {
                                ClusterID: clusterID,
                        },
                },
        })
        if err != nil {
                fmt.Printf("Error creating IRON client: %v\n", err)
                return
        }
        result, resp, err := client.Codes.CreateOrUpdateCode(iron.Code{
                Name:  "mytest",
                Image: "loafoe/siderite:latest",
        })
        fmt.Printf("%v %v %v\n", result, resp, err)
}

```

# Queueing a task
```go
package main

import (
        "fmt"
        "github.com/philips-software/go-hsdp-api/iron"
)

var (
    projectID = "yourIronProjectID"
    projectToken = "yourIronProjectToken"
    clusterID = "yourIronClusterID"
    taskName = ""
)

func main() {
        client, err := iron.NewClient(&iron.Config{
                ProjectID: projectID,
                Token:     projectToken,
                ClusterInfo: []iron.ClusterInfo{
                        {
                                ClusterID: clusterID,
                        },
                },
        })
        if err != nil {
                fmt.Printf("Error creating IRON client: %v\n", err)
                return
        }
        result, resp, err := client.Tasks.QueueTask(iron.Task{
                CodeName:  "mytask",
                Payload:   `{"foo": "bar"}`,
        })
        fmt.Printf("%v %v %v\n", result, resp, err)
}

```

# Encryption
Some Iron clusters expect the Payload of a task to be encrypted.
You can use the `iron.EncryptPayload` function for this.