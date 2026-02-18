// Package aws implements the AWS cloud adapter.
package aws

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vijayaxai/terraship/internal/cloud"
)

// Adapter implements cloud.Adapter for AWS
type Adapter struct {
	cfg       aws.Config
	ec2Client *ec2.Client
	s3Client  *s3.Client
	iamClient *iam.Client
	region    string
	profile   string
}

// NewAdapter creates a new AWS adapter
func NewAdapter() *Adapter {
	return &Adapter{}
}

// Name returns the provider name
func (a *Adapter) Name() cloud.Provider {
	return cloud.ProviderAWS
}

// Initialize sets up the AWS adapter
func (a *Adapter) Initialize(ctx context.Context, cloudConfig cloud.CloudConfig) error {
	var opts []func(*config.LoadOptions) error

	// Set region
	if cloudConfig.AWSRegion != "" {
		a.region = cloudConfig.AWSRegion
		opts = append(opts, config.WithRegion(cloudConfig.AWSRegion))
	} else if region := os.Getenv("AWS_REGION"); region != "" {
		a.region = region
	} else {
		a.region = "us-east-1" // default
		opts = append(opts, config.WithRegion(a.region))
	}

	// Set profile
	if cloudConfig.AWSProfile != "" {
		a.profile = cloudConfig.AWSProfile
		opts = append(opts, config.WithSharedConfigProfile(cloudConfig.AWSProfile))
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	a.cfg = cfg
	a.ec2Client = ec2.NewFromConfig(cfg)
	a.s3Client = s3.NewFromConfig(cfg)
	a.iamClient = iam.NewFromConfig(cfg)

	return nil
}

// DetectProvider attempts to detect if AWS is the provider
func (a *Adapter) DetectProvider(ctx context.Context) (bool, float64, error) {
	confidence := 0.0

	// Check for AWS environment variables
	if os.Getenv("AWS_REGION") != "" || os.Getenv("AWS_DEFAULT_REGION") != "" {
		confidence += 0.3
	}
	if os.Getenv("AWS_PROFILE") != "" {
		confidence += 0.2
	}
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		confidence += 0.3
	}

	// Check for AWS credentials file
	homeDir, _ := os.UserHomeDir()
	credFile := fmt.Sprintf("%s/.aws/credentials", homeDir)
	if _, err := os.Stat(credFile); err == nil {
		confidence += 0.2
	}

	return confidence > 0.5, confidence, nil
}

// ValidateCredentials checks if AWS credentials are valid
func (a *Adapter) ValidateCredentials(ctx context.Context) error {
	// Try to call STS GetCallerIdentity (lightweight call)
	// For now, we'll just check if we can list regions
	_, err := a.ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return fmt.Errorf("AWS credentials validation failed: %w", err)
	}

	return nil
}

// GetResourceStatus retrieves the current status of an AWS resource
func (a *Adapter) GetResourceStatus(ctx context.Context, resourceType, resourceID string) (*cloud.ResourceStatus, error) {
	status := &cloud.ResourceStatus{
		ResourceID:   resourceID,
		ResourceType: resourceType,
		Properties:   make(map[string]interface{}),
	}

	switch {
	case strings.HasPrefix(resourceType, "aws_instance"):
		return a.getEC2InstanceStatus(ctx, resourceID)
	case strings.HasPrefix(resourceType, "aws_s3_bucket"):
		return a.getS3BucketStatus(ctx, resourceID)
	case strings.HasPrefix(resourceType, "aws_iam_role"):
		return a.getIAMRoleStatus(ctx, resourceID)
	default:
		return status, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

func (a *Adapter) getEC2InstanceStatus(ctx context.Context, instanceID string) (*cloud.ResourceStatus, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := a.ec2Client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe EC2 instance: %w", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return &cloud.ResourceStatus{
			ResourceID:   instanceID,
			ResourceType: "aws_instance",
			Exists:       false,
		}, nil
	}

	instance := result.Reservations[0].Instances[0]
	tags := make(map[string]string)
	for _, tag := range instance.Tags {
		if tag.Key != nil && tag.Value != nil {
			tags[*tag.Key] = *tag.Value
		}
	}

	return &cloud.ResourceStatus{
		ResourceID:   instanceID,
		ResourceType: "aws_instance",
		Exists:       true,
		State:        string(instance.State.Name),
		Tags:         tags,
		Properties: map[string]interface{}{
			"instance_type":     instance.InstanceType,
			"availability_zone": *instance.Placement.AvailabilityZone,
			"public_ip":         instance.PublicIpAddress,
			"private_ip":        instance.PrivateIpAddress,
		},
	}, nil
}

func (a *Adapter) getS3BucketStatus(ctx context.Context, bucketName string) (*cloud.ResourceStatus, error) {
	// Check if bucket exists
	_, err := a.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return &cloud.ResourceStatus{
			ResourceID:   bucketName,
			ResourceType: "aws_s3_bucket",
			Exists:       false,
		}, nil
	}

	status := &cloud.ResourceStatus{
		ResourceID:   bucketName,
		ResourceType: "aws_s3_bucket",
		Exists:       true,
		Properties:   make(map[string]interface{}),
	}

	// Get bucket tagging
	tagsOutput, err := a.s3Client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		tags := make(map[string]string)
		for _, tag := range tagsOutput.TagSet {
			if tag.Key != nil && tag.Value != nil {
				tags[*tag.Key] = *tag.Value
			}
		}
		status.Tags = tags
	}

	// Get bucket encryption
	encryptionOutput, err := a.s3Client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil && encryptionOutput.ServerSideEncryptionConfiguration != nil {
		status.Properties["encryption_enabled"] = true
	} else {
		status.Properties["encryption_enabled"] = false
	}

	// Get bucket versioning
	versioningOutput, err := a.s3Client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		status.Properties["versioning_enabled"] = versioningOutput.Status == "Enabled"
	}

	return status, nil
}

func (a *Adapter) getIAMRoleStatus(ctx context.Context, roleName string) (*cloud.ResourceStatus, error) {
	input := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	result, err := a.iamClient.GetRole(ctx, input)
	if err != nil {
		return &cloud.ResourceStatus{
			ResourceID:   roleName,
			ResourceType: "aws_iam_role",
			Exists:       false,
		}, nil
	}

	tags := make(map[string]string)
	for _, tag := range result.Role.Tags {
		if tag.Key != nil && tag.Value != nil {
			tags[*tag.Key] = *tag.Value
		}
	}

	return &cloud.ResourceStatus{
		ResourceID:   roleName,
		ResourceType: "aws_iam_role",
		Exists:       true,
		Tags:         tags,
		Properties: map[string]interface{}{
			"arn":         *result.Role.Arn,
			"create_date": result.Role.CreateDate,
		},
	}, nil
}

// ValidateResourceCompliance checks resource compliance with policies
func (a *Adapter) ValidateResourceCompliance(ctx context.Context, resourceType string, resource map[string]interface{}, rules []cloud.ValidationRule) ([]cloud.ValidationResult, error) {
	// This is typically handled by the rules engine
	// Cloud-specific validations can be added here
	return []cloud.ValidationResult{}, nil
}

// DetectDrift compares planned state with actual cloud resources
func (a *Adapter) DetectDrift(ctx context.Context, plannedState map[string]interface{}, resourceType, resourceID string) (*cloud.ResourceStatus, error) {
	actualStatus, err := a.GetResourceStatus(ctx, resourceType, resourceID)
	if err != nil {
		return nil, err
	}

	if !actualStatus.Exists {
		actualStatus.DriftDetected = true
		actualStatus.DriftDetails = []string{"Resource does not exist in AWS"}
		return actualStatus, nil
	}

	// Compare planned vs actual - basic implementation
	// In production, this would be more sophisticated
	driftDetails := []string{}

	// Check tags
	if plannedTags, ok := plannedState["tags"].(map[string]interface{}); ok {
		for key, value := range plannedTags {
			if actualValue, exists := actualStatus.Tags[key]; !exists {
				driftDetails = append(driftDetails, fmt.Sprintf("Tag '%s' missing", key))
			} else if fmt.Sprint(value) != actualValue {
				driftDetails = append(driftDetails, fmt.Sprintf("Tag '%s' differs: planned=%v, actual=%v", key, value, actualValue))
			}
		}
	}

	if len(driftDetails) > 0 {
		actualStatus.DriftDetected = true
		actualStatus.DriftDetails = driftDetails
	}

	return actualStatus, nil
}

// ListResources lists AWS resources of a given type
func (a *Adapter) ListResources(ctx context.Context, resourceType string) ([]string, error) {
	switch {
	case strings.HasPrefix(resourceType, "aws_instance"):
		return a.listEC2Instances(ctx)
	case strings.HasPrefix(resourceType, "aws_s3_bucket"):
		return a.listS3Buckets(ctx)
	default:
		return nil, fmt.Errorf("listing not supported for resource type: %s", resourceType)
	}
}

func (a *Adapter) listEC2Instances(ctx context.Context) ([]string, error) {
	result, err := a.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	var instanceIDs []string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if instance.InstanceId != nil {
				instanceIDs = append(instanceIDs, *instance.InstanceId)
			}
		}
	}

	return instanceIDs, nil
}

func (a *Adapter) listS3Buckets(ctx context.Context) ([]string, error) {
	result, err := a.s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	var bucketNames []string
	for _, bucket := range result.Buckets {
		if bucket.Name != nil {
			bucketNames = append(bucketNames, *bucket.Name)
		}
	}

	return bucketNames, nil
}

// Close cleans up AWS adapter resources
func (a *Adapter) Close() error {
	// AWS SDK clients don't require explicit cleanup
	return nil
}
