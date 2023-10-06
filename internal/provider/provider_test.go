package provider

import (
	"os"
	"strings"
	"testing"

	"github.com/Think-iT-Labs/edc-connector-client-go/edc"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/assert"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the EDC client is properly configured.
	// It is also possible to use the EDC environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	//
	providerConfig = `

	provider "edc" {
		token = "1234"
		addresses = {
			default = "http://localhost:29193/api"
			management = "http://localhost:29193/api/v1/data"
			protocol = "http://localhost:29193/api/v1/ids"
			public = "http://localhost:29193/public"
			control = "http://localhost:29193/control"
		}
	}

`
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"edc": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func TestEDCProvider_Configure(t *testing.T) {
	type args struct {
		data EDCProviderModel
	}
	defaultAddress := "http://localhost:29193/api"
	managementAddress := "http://localhost:29193/api/v1/data"
	protocolAddress := "http://localhost:29193/api/v1/ids"
	publicAddress := "http://localhost:29193/public"
	controlAddress := "http://localhost:29193/control"
	authToken := "dummyToken"
	tests := []struct {
		name                           string
		args                           args
		env                            map[string]string
		expectedToken                  string
		expectedEdcAddresses           edc.Addresses
		diagnosticsPathErrorAttributes []string
		expectedError                  bool
	}{
		{
			name: "valid configuration",
			args: args{
				data: EDCProviderModel{},
			},
			env: map[string]string{
				"EDC_TOKEN":      authToken,
				"EDC_DEFAULT":    defaultAddress,
				"EDC_MANAGEMENT": managementAddress,
				"EDC_PROTOCOL":   protocolAddress,
				"EDC_PUBLIC":     publicAddress,
				"EDC_CONTROL":    controlAddress,
			},
			expectedToken: authToken,
			expectedEdcAddresses: edc.Addresses{
				Default:    &defaultAddress,
				Management: &managementAddress,
				Protocol:   &protocolAddress,
				Public:     &publicAddress,
				Control:    &controlAddress,
			},
			diagnosticsPathErrorAttributes: []string{},
			expectedError:                  false,
		},
		{
			name: "invalid configuration",
			args: args{
				data: EDCProviderModel{},
			},
			env: map[string]string{
				"EDC_TOKEN": authToken,
			},
			expectedToken:                  authToken,
			expectedEdcAddresses:           edc.Addresses{},
			diagnosticsPathErrorAttributes: []string{"Protocol", "Public", "Management", "Default", "Control"},
			expectedError:                  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldEnv := stashEnv()
			defer popEnv(oldEnv)

			for k, v := range tt.env {
				os.Setenv(k, v)
			}
			resp := provider.ConfigureResponse{}
			token, edcAddresses := validateProviderOptions(tt.args.data, &resp)

			assert.Equal(t, tt.expectedError, resp.Diagnostics.HasError())

			if tt.expectedError == false {
				assert.Equal(t, token, tt.expectedToken)
				assert.Equal(t, edcAddresses.Control, tt.expectedEdcAddresses.Control)
				assert.Equal(t, edcAddresses.Management, tt.expectedEdcAddresses.Management)
				assert.Equal(t, edcAddresses.Public, tt.expectedEdcAddresses.Public)
				assert.Equal(t, edcAddresses.Protocol, tt.expectedEdcAddresses.Protocol)
				assert.Equal(t, edcAddresses.Default, tt.expectedEdcAddresses.Default)
			} else {
				assert.ElementsMatch(t, tt.diagnosticsPathErrorAttributes, getErrorPaths(resp.Diagnostics.Errors()))
			}

		})
	}
}

func getErrorPaths(errors diag.Diagnostics) []string {
	var pathSteps []string
	for _, dd := range errors {
		diagWithPath, ok := dd.(diag.DiagnosticWithPath)
		if !ok {
			continue
		}
		pathSteps = append(pathSteps, diagWithPath.Path().Steps().String())
	}
	return pathSteps
}

func stashEnv() []string {
	env := os.Environ()
	os.Clearenv()
	return env
}

func popEnv(env []string) {
	os.Clearenv()

	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		k, v := p[0], ""
		if len(p) > 1 {
			v = p[1]
		}
		os.Setenv(k, v)
	}
}
