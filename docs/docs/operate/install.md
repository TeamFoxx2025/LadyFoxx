!!! tip

    We recommend using the pre-built releases and verifying the provided checksums for security.

    The Docker image is also a convenient option for containerized deployment. Building from source provides greater flexibility, but requires a [<ins>suitable development environment</iedgens>](system.md).

## Pre-built releases

To access the pre-built releases, visit the [<ins>GitHub releases page</ins>](https://github.com/TeamFoxx2025/LadyFoxx/releases). 
The client provides cross-compiled AMD64/ARM64 binaries for Darwin and Linux.

!!! info "Latest release: 1.3.0"

    **The latest stable test release is [<ins>v1.3.0</ins>](https://github.com/TeamFoxx2025/LadyFoxx/releases/tag/v1.3.0).**

## Docker image

To use Docker, you will need to have it installed on your system. If you haven't already installed Docker, you can follow the instructions on the
[<ins>official Docker website</ins>](https://www.docker.com/) for your operating system.

You can access the official LadyFoxx Docker images hosted under the [<ins> registry</ins>](https://hub.docker.com/r/TeamFoxx2025/LadyFoxx) using the following command:

  ```bash
  docker pull TeamFoxx2025/LadyFoxx:latest
  ```

## Build from source

> Before getting started, ensure you have [Go](https://go.dev/) installed on your system (version >= 1.15 and <= 1.19).
> Compatibility is being worked on for other versions and will be available in the near future.

> To install Go, run the following command in your CLI (we are using 1.18 in this example): `sudo apt-get install golang-1.18`.
> Or, use a package manager like [<ins>Snapcraft</ins>](https://snapcraft.io/go) for Linux, [<ins>Homebrew</ins>](https://formulae.brew.sh/formula/go) for Mac, and [<ins>Chocolatey</ins>](https://community.chocolatey.org/packages/golang) for Windows.

Use the following commands to clone the LadyFoxx repository and build from source:

  ```bash
  git clone https://github.com/TeamFoxx2025/LadyFoxx.git
  cd foxx-chain/
  go build -o foxx-chain .
  ```
