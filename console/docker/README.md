# HSP Docker

[![Go Reference](https://pkg.go.dev/badge/github.com/philips-software/go-hsdp-api@main/console/docker.svg)](https://pkg.go.dev/github.com/philips-software/go-hsdp-api@main/console/docker)

The docker client provides access to the HSDP Docker registry.

## Example usage

The following example will create a `console.Client` and use it to create a `docker.Client` to list all namespaces and repositories in the `us-east` region:

```golang
package main

import (
	"context"
	"fmt"

	"github.com/philips-software/go-hsdp-api/console"
	"github.com/philips-software/go-hsdp-api/console/docker"
)

func main() {
	consoleClient, err := console.NewClient(nil, &console.Config{
		Region:   "us-east",
	})
	if err != nil {
		fmt.Printf("Error creating console client: %v\n", err)
	}
	err = consoleClient.Login("cfusernamehere", `cfpasswordhere`)
	if err != nil {
		fmt.Printf("Error logging into Console: %v\n", err)
		return
	}
	dockerClient, err := docker.NewClient(consoleClient, &docker.Config{
		Region:   "us-east",
	})
	if err != nil {
		fmt.Printf("Error creating docker client: %v\n", err)
	}
	ctx := context.Background()
	list, err := dockerClient.Namespaces.GetNamespaces(ctx)

	if err != nil {
		fmt.Printf("Error getting namespaces: %v\n", err)
		return
	}
	for _, namespace := range *list {
		fmt.Printf("------ Namespace: %s --------\n", namespace.ID)
		repos, err := dockerClient.Namespaces.GetRepositories(ctx, namespace.ID)
		if err != nil {
			fmt.Printf("Error fetching repo '%s': %v\n", namespace.ID, err)
			continue
		}
		for _, repo := range *repos {
			latest, err := dockerClient.Repositories.GetLatestTag(ctx, repo.ID)
			tag := ""
			digest := ""
			if err == nil {
				tag = ":" + latest.Name

			}
			fmt.Printf("Repo: %s/%s%s (%s)\n", namespace.ID, repo.Name, tag, digest)
		}
	}
}
```
