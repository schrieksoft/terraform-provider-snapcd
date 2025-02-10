// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"fmt"
	"math/rand"
	provider "terraform-provider-snapcd/internal/provider"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "snapcd" {
	url                           = "https://localhost:20002"
	insecure_skip_verify          = true
	health_check_interval_seconds = 5
	health_check_timeout_seconds  = 20
	client_id 					  = "Terraformer"
	client_secret 				  = "somesecretforsampleclient"
	# access_token 				  = "eyJhbGciOiJBMjU2S1ciLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwidHlwIjoiYXQrand0IiwiY3R5IjoiSldUIn0.ipCqtoL90B4ZshniO3VqJZ4eJpzZfrc0TotbHS77L6URvLCxBXAIX_i5XeOEb53z35dIk5RX-HHkX1MPjcR2lL21KgXnwo2Z.he8m3gEnsfOSn4gaqFNuzQ._d4b8JAVO_tUoSdGB71uCn4ng-mo_vMuQzpXvKXvKUdjl6McEEWzVYCJgUd_VX87D7h5UPlRGgrzSluSLUltITOzJVchzbU0hOJ8Kv-Pt2dnFhS8EhhAg0kEs2VSmmHkKMUJdudOfGOQsp1IeZagObhDg-Xorv4sSSiQLry-blq0KYPEqRAcCid_tK3YEQiaOEEwqdQVbDzq16VbBDdTFtUV1IjYxkWt1EQoigelBNuieaGJs84DlBXBnulede7KJgU7Z_2AulsoTgjxC-PkRr2cLYaKZzH836p2YInI609Fh76wdgb5oo95gXpKmVCgi_QufdpgGC436qK0qkeQ4wBJ9_A49hmzlcTOGIWoY1uUuvbjoFltcvyqxCTctbd0E4ATVn7JidTqlWNa__zZYnSNyvXx0iDEZLlp8whdTETtOKzcBFEeoMGIMHToaIUjBD_YyFZzQtjb_0KRStq5MNTWFCbjvtQ8s0nysxJ_L2g0UKL_q4GjoCBwZI4rhSyAl1ZBvvnnFZtp8g0Hi5SjxhH0GXqAmcvs51fFKwT3wqobA_BCqxqgwLEfvn8rlT9gz44QKWrzRitkmfRvwPviSUoQmiMySrLInbmZ5x50ydAPm2i_bE3JEEyjkzt_j3hJ-8AwxZRE-9NTkSJABexL7WUVs585judvL8Edu_BFY-YA6gBK6piKynjTuSFrkAYq1pyMixJZLFvYtGQJJlkRxxD44q2DI7Yr04oJ-6Krg-FQZU92q3AcCoywrBh7E-gnVVbtrmyC-LpaqLzJxJmXRfOdhUoWmCSCYnEYnU2VHgo4pnWjWQqfs_KHGoSQEW1576nLgjkMupkrK-CgNimtFkkd-W0g3znp8Ypf2nIiUerG3HGpH4CiOBqeHEwT2ZHsCAM1tnUKbgfSmFP8rGakEpzbaZcpFCkWUkgPzXhOx2urVi727lVmxV1DGqfWxMMNDudFn1FQ9Y1l5DKKOZnJ8-TqDK86diTtQcCL1GZQ2nPvH3Z9Ds6wSH_zfAmqPhMW0ZrcTq3DUwp4MmQN5EWHQY3CdlqcIK8F_PujenXQ6f-lY_6VsjoQSyaPnMq39WlrBM-O-ZL4CxAG4V_fyOLcVABMDRsoqH1HKofOepen45Q.xYw7jXkFi2D1TQ8uDR1mpW41I9neBCEwkv2_q8AGK_Y"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"snapcd": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
	namePostfix = generateRandomString(16)
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

func appendRandomString(input string) string {
	return fmt.Sprintf(input, namePostfix)
}
