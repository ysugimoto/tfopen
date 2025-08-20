# tfopen

[![build](https://github.com/ysugimoto/tfopen/actions/workflows/build.yml/badge.svg)](https://github.com/ysugimoto/tfopen/actions/workflows/build.yml)
[![release](https://github.com/ysugimoto/tfopen/actions/workflows/release.yml/badge.svg)](https://github.com/ysugimoto/tfopen/actions/workflows/release.yml)

`tfopen` is a simple CLI tool that reads your Terraform configuration files (`.tf`) and opens the corresponding Terraform Cloud/Enterprise workspace UI in your browser.

It scans `.tf` files in the current directory, finds the workspace information from either a `cloud {}` block or a `backend "remote" {}` block, and constructs the workspace URL to open.

## Features

- **Auto Discovery**: Automatically scans `.tf` files in the current directory.
- **Multiple Block Support**: Supports both `cloud {}` and `backend "remote" {}` blocks.
- **Cross-Platform**: Works on macOS, Linux, and Windows.

## Installation

### Using Go

```bash
go install github.com/ysugimoto/tfopen@latest
```

### Using Homebrew (macOS and Linux)

```bash
brew install ysugimoto/tap/tfopen
```

### Manual Installation

Download the appropriate binary for your OS from the [GitHub Releases](https://github.com/ysugimoto/tfopen/releases) page.

## Usage

Simply run the following command in the directory where your Terraform configuration files are located:

```bash
tfopen
```

The tool will parse the `.tf` files and automatically open the Terraform Cloud workspace page in your browser.

### Supported Terraform Configurations

`tfopen` supports Terraform configurations written in either of the following formats:

#### `cloud {}` block

If a `cloud {}` block is defined within the `terraform` block, `tfopen` will use it.

**Example:** `main.tf`
```hcl
terraform {
  cloud {
    organization = "my-organization"

    workspaces {
      name = "my-workspace"
    }
  }
}
```

#### `backend "remote" {}` block

If a `backend` block is defined with the `"remote"` type, `tfopen` will use it.

**Example:** `main.tf`
```hcl
terraform {
  backend "remote" {
    organization = "my-organization"

    workspaces {
      name = "my-workspace"
    }
  }
}
```

## License

[MIT License](./LICENSE)