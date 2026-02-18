package terraform

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tmpDir := t.TempDir()

	client, err := NewClient(tmpDir)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, tmpDir, client.workingDir)
}

func TestNewClient_InvalidDir(t *testing.T) {
	client, err := NewClient("/nonexistent/directory")
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestClient_SetEnvironment(t *testing.T) {
	tmpDir := t.TempDir()
	client, _ := NewClient(tmpDir)

	client.SetEnvironment("AWS_REGION", "us-west-2")
	assert.Equal(t, "us-west-2", client.envVars["AWS_REGION"])
}

func TestClient_GetProvider(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test Terraform file
	tfContent := `
provider "aws" {
  region = "us-east-1"
}

resource "aws_instance" "example" {
  ami           = "ami-12345"
  instance_type = "t2.micro"
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.tf"), []byte(tfContent), 0644)
	require.NoError(t, err)

	client, _ := NewClient(tmpDir)
	ctx := context.Background()

	provider, err := client.GetProvider(ctx)
	require.NoError(t, err)
	assert.Equal(t, "aws", provider)
}
