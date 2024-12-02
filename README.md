# Cloudflare Tool (cftool)

## Overview

The Cloudflare Tool (cftool) is a command-line utility written in Go that interacts with the Cloudflare API to manage various settings for a specified domain. It allows users to perform actions such as purging the cache, toggling development mode, and adjusting the security level of a Cloudflare zone.

## Features

- **Purge Cache**: Clear the cache for a specified Cloudflare zone.
- **Development Mode**: Enable or disable development mode for a specified zone.
- **Security Level**: Change the security level of a specified zone to one of the following: `attack`, `high`, `medium`, or `low`.

## Usage

### Flags

- `-email`: Cloudflare account email. This is required unless set via the `CF_EMAIL` environment variable.
- `-key`: Cloudflare API key. This is required unless set via the `CF_KEY` environment variable.
- `-zone`: Cloudflare zone (domain) to manage. This is required unless set via the `CF_DOMAIN` environment variable.
- `-purge-cache`: Boolean flag to purge the cache for the specified zone.
- `-development-mode`: Set to `on` or `off` to enable or disable development mode.
- `-secure-level`: Set the security level to `attack`, `high`, `medium`, or `low`.

### Example Commands

- Purge Cache:
  ```bash
  cftool -email your-email@example.com -key your-api-key -zone example.com -purge-cache
  ```

- Enable Development Mode:
  ```bash
  cftool -email your-email@example.com -key your-api-key -zone example.com -development-mode on
  ```

- Set Security Level to High:
  ```bash
  cftool -email your-email@example.com -key your-api-key -zone example.com -secure-level high
  ```

## Error Handling

The program will output error messages if any required flags are missing or if an invalid value is provided for a flag. It will also display usage information if no actions are specified.

## Dependencies

- [Go](https://golang.org/)
- [tidwall/gjson](https://github.com/tidwall/gjson) for JSON parsing

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/cftool.git
   ```

2. Navigate to the project directory:
   ```bash
   cd cftool
   ```

3. Get requires:
   ```bash
   go mod tidy
   ```

4. Build the project:
   ```bash
   go build -o cftool
   ```

5. Run the tool using the examples provided above.

## License

This project is licensed under the Apache License Version 2.0.