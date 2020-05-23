terraform {
  required_version = "{{.TerraformVersion}}"

  required_providers {
    "{{.Provider}}" = "{{.ProviderVersion}}"
  }
}
