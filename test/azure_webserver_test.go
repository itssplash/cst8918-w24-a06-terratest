package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "454cef6e-c0dd-444f-9766-2fcc4e902f67"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "abhi0012",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Confirm NIC exists and is connected to VM
	nics := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)
	assert.NotEmpty(t, nics, "Expected at least one NIC to be attached to the VM")

	// Confirm the VM is running the correct Ubuntu version
	vmImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)
	// Check publisher and offer for Ubuntu
	assert.Equal(t, "Canonical", vmImage.Publisher, "Expected publisher to be Canonical")
	assert.Equal(t, "0001-com-ubuntu-server-jammy", vmImage.Offer, "Expected offer to be 0001-com-ubuntu-server-jammy")
	assert.Equal(t, "22_04-lts-gen2", vmImage.SKU, "Expected SKU to be 22_04-lts-gen2") //checking if the same version of ubuntu is installed on VM
}
