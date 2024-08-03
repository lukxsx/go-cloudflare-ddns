# go-cloudflare-ddns

Simple Dynamic DNS client for Cloudflare-managed domains written in the Go programming language.

## Configuration

The client is configured using environment variables. The client can read the variables from `.env` file.

| **Variable**      | **Explanation**                                                              | **Default** |
|-------------------|------------------------------------------------------------------------------|-------------|
| `CF_API_TOKEN`    | Cloudflare API token                                                         | none        |
| `CF_ZONE_ID`      | Cloudflare Zone ID for the DNS zone to manage                                | none        |
| `DOMAINS`         | Comma-separated list of domains that the client controls                     | none        |
| `DEBUG`           | Set to `true` to enable debug printing. Warning: will print unmasked secrets | `false`     |
| `UPDATE_INTERVAL` | Frequency at which the client checks for changes. The value is in minutes.   | `10`        |
