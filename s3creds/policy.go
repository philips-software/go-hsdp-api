package s3creds

import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Policy struct {
	ID           int    `json:"id,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
	Conditions   struct {
		ManagingOrganizations []string `json:"managingOrganizations,omitempty"`
		Groups                []string `json:"groups,omitempty"`
	} `json:"conditions,omitempty"`
	Allowed struct {
		Resources []string `json:"resources" validate:"min=1"`
		Actions   []string `json:"actions" validate:"policyActions,min=1,unique"`
	} `json:"allowed"`

	ProductKey string `json:"-"`
}

func (p *Policy) StringID() string {
	return strconv.Itoa(p.ID)
}

// Equals determines of other Policy is equavalent
func (p *Policy) Equals(other *Policy) bool {
	if p.ID != other.ID {
		return false
	}
	if p.ResourceType != other.ResourceType {
		return false
	}
	if len(difference(p.Conditions.ManagingOrganizations, other.Conditions.ManagingOrganizations)) > 0 {
		return false
	}
	if len(difference(other.Conditions.ManagingOrganizations, p.Conditions.ManagingOrganizations)) > 0 {
		return false
	}
	if len(difference(p.Conditions.Groups, other.Conditions.Groups)) > 0 {
		return false
	}
	if len(difference(other.Conditions.Groups, p.Conditions.Groups)) > 0 {
		return false
	}
	if len(difference(p.Allowed.Resources, other.Allowed.Resources)) > 0 {
		return false
	}
	if len(difference(other.Allowed.Resources, p.Allowed.Resources)) > 0 {
		return false
	}
	if len(difference(p.Allowed.Actions, other.Allowed.Actions)) > 0 {
		return false
	}
	if len(difference(other.Allowed.Actions, p.Allowed.Actions)) > 0 {
		return false
	}
	if p.ProductKey != other.ProductKey {
		return false
	}
	return true
}

func validateActions(fl validator.FieldLevel) bool {
	validActions := []string{
		"GET",
		"PUT",
		"LIST",
		"DELETE",
		"ALL_OBJECT",
		"ALL_BUCKET",
		"AbortMultipartUpload",
		"CompleteMultipartUpload",
		"CopyObject",
		"CreateBucket",
		"CreateMultipartUpload",
		"DeleteBucket",
		"DeleteBucketAnalyticsConfiguration",
		"DeleteBucketCors",
		"DeleteBucketEncryption",
		"DeleteBucketIntelligentTieringConfiguration",
		"DeleteBucketInventoryConfiguration",
		"DeleteBucketLifecycle",
		"DeleteBucketMetricsConfiguration",
		"DeleteBucketOwnershipControls",
		"DeleteBucketPolicy",
		"DeleteBucketReplication",
		"DeleteBucketTagging",
		"DeleteBucketWebsite",
		"DeleteObject",
		"DeleteObjects",
		"DeleteObjectTagging",
		"DeletePublicAccessBlock",
		"GetBucketAccelerateConfiguration",
		"GetBucketAcl",
		"GetBucketAnalyticsConfiguration",
		"GetBucketCors",
		"GetBucketEncryption",
		"GetBucketIntelligentTieringConfiguration",
		"GetBucketInventoryConfiguration",
		"GetBucketLifecycle",
		"GetBucketLifecycleConfiguration",
		"GetBucketLocation",
		"GetBucketLogging",
		"GetBucketMetricsConfiguration",
		"GetBucketNotification",
		"GetBucketNotificationConfiguration",
		"GetBucketOwnershipControls",
		"GetBucketPolicy",
		"GetBucketPolicyStatus",
		"GetBucketReplication",
		"GetBucketRequestPayment",
		"GetBucketTagging",
		"GetBucketVersioning",
		"GetBucketWebsite",
		"GetObject",
		"GetObjectAcl",
		"GetObjectLegalHold",
		"GetObjectLockConfiguration",
		"GetObjectRetention",
		"GetObjectTagging",
		"GetObjectTorrent",
		"GetPublicAccessBlock",
		"HeadBucket",
		"HeadObject",
		"ListBucketAnalyticsConfigurations",
		"ListBucketIntelligentTieringConfigurations",
		"ListBucketInventoryConfigurations",
		"ListBucketMetricsConfigurations",
		"ListBuckets",
		"ListMultipartUploads",
		"ListObjects",
		"ListObjectsV2",
		"ListObjectVersions",
		"ListParts",
		"PutBucketAccelerateConfiguration",
		"PutBucketAcl",
		"PutBucketAnalyticsConfiguration",
		"PutBucketCors",
		"PutBucketEncryption",
		"PutBucketIntelligentTieringConfiguration",
		"PutBucketInventoryConfiguration",
		"PutBucketLifecycle",
		"PutBucketLifecycleConfiguration",
		"PutBucketLogging",
		"PutBucketMetricsConfiguration",
		"PutBucketNotification",
		"PutBucketNotificationConfiguration",
		"PutBucketOwnershipControls",
		"PutBucketPolicy",
		"PutBucketReplication",
		"PutBucketRequestPayment",
		"PutBucketTagging",
		"PutBucketVersioning",
		"PutBucketWebsite",
		"PutObject",
		"PutObjectAcl",
		"PutObjectLegalHold",
		"PutObjectLockConfiguration",
		"PutObjectRetention",
		"PutObjectTagging",
		"PutPublicAccessBlock",
		"RestoreObject",
		"SelectObjectContent",
		"UploadPart",
		"UploadPartCopy",
		"WriteGetObjectResponse",
	}
	// Also need to check for granular objects
	actions, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}
	for _, a := range actions {
		found := false
		for _, v := range validActions {
			if a == v {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
