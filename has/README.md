# Using the HAS API client
Hosted Application Streaming is a HSDP service which provides application deployment in the cloud or on premise with streaming output to your browser. This allows one to quickly scale up and consume workstation class hardware on any browser enabled system.

# Provisioning HAS resources

```golang
import (
	"fmt"

	"github.com/philips-software/go-hsdp-api/has"
	"github.com/philips-software/go-hsdp-api/iam"
)

func main() {
	iamClient, err := iam.NewClient(nil, &iam.Config{
		OAuth2ClientID: "YourClientID",
		OAuth2Secret:   "YourClientSecret",
		Region: "eu-west",
		Environment: "client-test",
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
	}
	err = iamClient.Login("yourlogin", "yourpassword")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	client, err := has.NewClient(iamClient, &has.Config{
		HASURL:   "https://has-client-test.eu-west.philips-healthsuite.com",
		OrgID:    "your-org-uuid-here",
	})
	if err != nil {
		fmt.Printf("Error creating HAS client: %v\n", err)
		return
	}
	images, resp, err := client.Images.GetImages()
	fmt.Printf("%v %v %v", images, resp, err)

	res, resp, err := client.Resources.CreateResource(has.Resource{
		ImageID:      "ami-0fc5fakeimageidhere",
		ResourceType: "g3s.xlarge",
		Region:       "eu-west-1",
		Count:        1,
		ClusterTag:   "andytest",
		EBS: has.EBS{
                        DeleteOnTermination: true,
			VolumeSize: 50,
			VolumeType: "standard",
		},
	})

	fmt.Printf("%v %v %v", res, resp, err)
}
```
