# GPX (Golang Package eXecute)

GPX is a command-line tool inspired by `npx`, designed to execute Go packages
directly from your `go.mod` file. It ensures that you always run the version
specified in your project, making it easier to manage dependencies and execute
tools without worrying about version mismatches.

## Installation

To install GPX, you can use the following command:

```bash
go install github.com/Polo123456789/gpx@latest
```

## Usage

To use GPX, simply run the following command in your terminal:

```bash
gpx <tool> [args...]
```

## Embeded Commands

* `i:install` Installs all tools listed in the `tools.go` file.
* `i:clean` Removes all tools that you haven't used in the last 30 days.

## Setting Up Your Tools

To set up your tools, create a `tools.go` file in your project directory.
Here’s an example of how to structure it:

```go
//go:build tools
// +build tools

package tools

import (
    _ "github.com/some/tool" // Replace with your tool's import path
    _ "github.com/another/tool"
)
```

Make sure to run `go mod tidy` after adding tools to ensure they are included
in your `go.mod` file.

## Contributing

Contributions are welcome! If you have suggestions or improvements, feel free
to open an issue or submit a pull request.

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/YourFeature`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a pull request.
