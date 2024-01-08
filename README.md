# shared-broker-utils

This repo contains a Go module that provides helpers for the cloud.gov service brokers written in Go, including:

- <https://github.com/cloud-gov/aws-broker>
- <https://github.com/cloud-gov/s3-broker>
- <https://github.com/cloud-gov/uaa-credentials-broker>

The helpers included in this module include:

- Helper function for generating tags for provisioned resources. Based on the GUIDs provided, the function also uses the CF API to look up names of the associated resources
