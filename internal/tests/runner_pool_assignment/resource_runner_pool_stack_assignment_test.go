// SPDX-License-Identifier: MPL-2.0

package runner_pool_assignment

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRunnerPoolStackAssignment_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + core.RunnerPoolCreateConfig + RunnerPoolStackAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool_stack_assignment.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceRunnerPoolStackAssignment_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + core.RunnerPoolCreateConfig + RunnerPoolStackAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool_stack_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_runner_pool_stack_assignment.this", "permission", "ReadWrite"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_runner_pool_stack_assignment" "this" { 
  runner_pool_id   = snapcd_runner_pool_stack_assignment.this.id
  stack_id   = snapcd_stack.this.id
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool_stack_assignment.this", "id"),
					resource.TestCheckResourceAttr("snapcd_runner_pool_stack_assignment.this", "permission", "Read"),
				),
			},
		},
	})
}

func TestAccResourceRunnerPoolStackAssignment_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.StackCreateConfig + core.RunnerPoolCreateConfig + RunnerPoolStackAssignmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_runner_pool_stack_assignment.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_runner_pool_stack_assignment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
