# gosect

A lightweight tool to dynamically replace content sections in files. Keep your
documentation, code snippets, and configuration files always up-to-date by
defining sections that automatically sync with source files.

Perfect for maintaining README files with live code examples, keeping
documentation in sync with actual configuration, or any scenario where you need
to embed external content into your files.

## Features

- **Zero Dependencies**: Single binary, ready to use
- **Custom Markers**: Support for any comment style (HTML, Markdown, code
  comments, etc.)
- **Multiple Sections**: Handle multiple sections in a single file
- **File References**: Each section can reference a different source file
- **Preserves Formatting**: Keeps your original file structure intact
- **Docker Ready**: Multi-arch Docker images (amd64/arm64)
- **Lightweight**: Fast and efficient

## Use Cases

- **Live Code Examples**: Keep README code examples in sync with actual working
  files
- **Configuration Documentation**: Embed real config snippets in your docs
- **Multi-file Documentation**: Compose documentation from multiple source files
- **Template Management**: Maintain templates with dynamic content sections
- **CI/CD Integration**: Automatically update documentation in your pipeline

## Quick Start

### Standalone Binary

```bash
# Install
go install github.com/badele/gosect@latest

# Mark sections in your file
cat > README.md << 'EOF'
# My Project

## Installation

<!-- BEGIN SECTION install file=./install.sh -->
<!-- END SECTION install -->

## Configuration

<!-- BEGIN SECTION config file=./config.yaml -->
<!-- END SECTION config -->
EOF

# Create source files
echo "npm install my-package" > install.sh
echo "port: 8080\nhost: localhost" > config.yaml

# Replace sections
gosect -file README.md

# Your README.md now contains the actual content!
```

### Docker

```bash
# Run with Docker
docker run -v .:/work badele/gosect:latest -file /work/README.md
```

### Nix

```bash
# Run directly with Nix
nix run github:badele/gosect -- -file README.md

# Or add to your flake.nix
```

## Usage

### Command Line Options

```
Usage of gosect:
  -file string
        Input file path (required)
  -begin string
        Begin marker prefix (default "BEGIN SECTION")
  -end string
        End marker prefix (default "END SECTION")
  -stdout
        Print to stdout instead of writing file
  -verbose
        Log details about processed sections
```

### Section Syntax

Mark sections in your files using this format:

```
<marker-prefix> BEGIN SECTION <section-name> file=<source-file>
content will be replaced here
<marker-prefix> END SECTION <section-name>
```

#### Examples

```markdown
<!-- BEGIN SECTION install file=./install.sh -->
<!-- END SECTION install -->
```

```python
# BEGIN SECTION config file=./config.json
# END SECTION config
```

```html
<!-- BEGIN SECTION footer file=./footer.html -->
<!-- END SECTION footer -->
```

#### Example 2: Custom Markers

You can use any marker style that suits your file format:

```bash
# For shell scripts
gosect -file script.sh -begin "# START" -end "# STOP"

# For C/C++ comments
gosect -file code.cpp -begin "// BEGIN" -end "// END"

# For configuration files
gosect -file config.ini -begin "; BEGIN" -end "; END"
```
