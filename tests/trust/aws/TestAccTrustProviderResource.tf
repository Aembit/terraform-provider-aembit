provider "aembit" {
}

resource "aembit_trust_provider" "aws" {
	name = "TF Acceptance AWS"
	is_active = true
	aws_metadata = {
		certificate = <<-EOT
-----BEGIN CERTIFICATE-----
MIIDgTCCAmmgAwIBAgIQSpi05faO6r1DOItzHn1SfjANBgkqhkiG9w0BAQsFADAX
MRUwEwYDVQQDDAxhZW1iaXQubG9jYWwwHhcNMjQwNTMxMTM0MjA5WhcNMjUwNTMx
MTQwMjA5WjAXMRUwEwYDVQQDDAxhZW1iaXQubG9jYWwwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDgC/cb8YsJExhNgu1mZERrU8RrOWu6qtV/jfuN4lgJ
WvRuSq8uLkN6ZQA36s42aNy2hy6smxJZZWuQfiWXeLVPICcuXcITlMbpPnok/0LS
0o3jLQcw/rT/i2ybdRanko1/Vag8RZOGaIDHLyfeuoc1RQ2yZB6hC+8eUkZEL3XH
kTwQ+uD/blcpgGYFN7B6bkhpUB7dA9xf5JOOA7hQVG9RNt3bsQxyZndGBPKuac27
rL+2RhDQuns+ahasRRjQyV4uvgGNM4X6aoylVnct7KN/PV/iKX4RkK3cLRO4CnuB
nksTE13NZWi4LyPNTDMO8odeSwC3yfUHkEwNDiRACToJAgMBAAGjgcgwgcUwDgYD
VR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcDATB1BgNV
HREEbjBsggxhZW1iaXQubG9jYWyCDiouYWVtYml0LmxvY2FsghIqLmFwcC5hZW1i
aXQubG9jYWyCEiouYXBpLmFlbWJpdC5sb2NhbIIRKi5pZC5hZW1iaXQubG9jYWyC
ESouZWMuYWVtYml0LmxvY2FsMB0GA1UdDgQWBBRwId7rKwY+2QkkdcWQkQKHzSYu
CTANBgkqhkiG9w0BAQsFAAOCAQEAzFc6UQdoUuit5bQxvadh/aZWnkOe/aL+q/vo
B+ONcymqB7QAXOJFDbUm/4Fq4xb8yj/Y+vF6Vbe5IM/esIRkKzLkT/Fa+yWmGx66
pJRhOrS8LYVJMnuJ3sInfAT8aLIYfKsYAkyMnJ70oWRmWbVvo5L+ukqmFk8WzWUi
aXC/amSYD82+BouCdHiufbFADS0WbpJTN25hJhPlRNaIRxzD63m3zjWQtEK+oUwD
uy16zlW2rUwNgW/Bgl8RLuf3gZnUStBoLV3xqIA3BibEp4DpYLj4AKz16LM/b/33
sfa3cNOSjZ3VJtS4B9JoBHh6RrAIxSk2Hoz9dTTDrp61oz3Lrw==
-----END CERTIFICATE-----
EOT
		account_id = "account_id"
		architecture = "architecture"
		availability_zone = "availability_zone"
		billing_products = "billing_products"
		image_id = "image_id"
		instance_id = "instance_id"
		instance_type = "instance_type"
		kernel_id = "kernel_id"
		marketplace_product_codes = "marketplace_product_codes"
		pending_time = "pending_time"
		private_ip = "private_ip"
		ramdisk_id = "ramdisk_id"
		region = "region"
		version = "version"
	}
}