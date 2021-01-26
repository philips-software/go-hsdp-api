# Using the Cartel API client
Container Host is a HSDP service which provides a hardened Docker runtime environment. This API client allows one to spin up Container Host instances using a convenient API. Some examples below:

# Spinning up a Container Host instance

```golang
package main

import (
	"github.com/philips-software/go-hsdp-api/cartel"
)

func main() {
	client, err := cartel.NewClient(nil, cartel.Config{
		Token:  "YourCartelToken",
		Secret: "YourCartelSecr3t",
		Host:   "cartel-host.here.com",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	myinstance, _, err := client.Create("myinstance.dev",
		cartel.EncryptVolumes(),
		cartel.VolumesAndSize(1, 50),
		cartel.SecurityGroups("https-from-cf", "tcp-1080"),
		cartel.UserGroups("my-ldap-group"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("InstanceID: %s\n", myinstance.InstanceID())
	fmt.Printf("InstanceIP: %s\n", myinstance.IPAddress())
}
```

# Get instance details

```golang
package main

import (
	"github.com/philips-software/go-hsdp-api/cartel"
)

func pretty(data []byte) string {
	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, data, "", "    ")
	return string(prettyJSON.Bytes())
}

func main() {
	client, err := cartel.NewClient(nil, cartel.Config{
		Token:  "YourCartelToken",
		Secret: []byte("YourCartelSecr3t"),
		Host:   "cartel-host.here.com",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	details, _, err := client.Details("myinstancer.dev")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%v\n", pretty(details))
}
```

# Destroying an instance

```golang
package main

import (
	"github.com/philips-software/go-hsdp-api/cartel"
)

func main() {
	client, err := cartel.NewClient(nil, cartel.Config{
		Token:  "YourCartelToken",
		Secret: []byte("YourCartelSecr3t"),
		Host:   "cartel-host.here.com",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	result, _, err := client.Destroy("myinstancer.dev")

	fmt.Printf("Result: %v\n", result.Success())
}
```
