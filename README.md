# go-cloudflare-ddns

Simple Dynamic DNS client for Cloudflare-managed domains, written in the Go programming language.

Currently, only IPv4 IPs and A DNS records are supported. IPv6 support is planned.

## Configuration

The client is configured using environment variables. The client can read the variables from `.env` file.

| **Variable**      | **Explanation**                                                              | **Default** |
|-------------------|------------------------------------------------------------------------------|-------------|
| `CF_API_TOKEN`    | Cloudflare API token                                                         | none        |
| `CF_ZONE_ID`      | Cloudflare Zone ID for the DNS zone to manage                                | none        |
| `DOMAINS`         | Comma-separated list of domains that the client controls                     | none        |
| `DEBUG`           | Set to `true` to enable debug printing. Warning: will print unmasked secrets | `false`     |
| `UPDATE_INTERVAL` | Frequency at which the client checks for changes. The value is in minutes.   | `10`        |

## Usage

You should have a Cloudflare zone created for your domain. The zone should have A records set to match your domains.

1. Create a API token in the [Cloudflare dashboard](https://dash.cloudflare.com/profile/api-tokens).

2. Set the API token and the Zone ID as the `CF_API_TOKEN` and `CF_ZONE_ID` environment variables.

3. Set the domains to change in the `DOMAINS` variable as an comma-separated list. Example: `example.com,subdomain.example.com`

4. Run the application using `go run .`


## Planned features

* IPv6 support
* Automatic creation of missing DNS records
