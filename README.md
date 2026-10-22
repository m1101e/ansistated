<div align="right">
  <a href="README_RU.md">Русский</a> | <b>English</b>
</div>

# ansistated

**ansistated** is a lightweight CLI templating engine and variable manager tailored for dotfiles and configuration files.

> Originally developed for the [ansistate](https://github.com/m1101e/ansistate) environment configuration, this tool was designed with versatility in mind and can be utilized with any templates and environments.

---

## Features

- **Pongo2-based Templating (Django/Jinja2 syntax)** with built-in filters for seamless color format conversions (`to_hex`, `to_rgb`, `to_strip_hex`).
- **YAML Frontmatter in Templates**: Each template self-defines its target destination path (`dest`), target groups (`group`), and optional post-render execution hook (`post_cmd`).
- **Precise YAML Updates (AST-based)**: The `update` command updates variables from a JSON stream while preserving comments, indentation, and structure of your YAML files.
- **Concurrent Rendering**: High-performance worker-pool processing with file-content diffing to avoid unnecessary writes to disk.

---

## Installation

### Option 1. Pre-built Binaries (Recommended)

Download the archive for your operating system and CPU architecture from the [Releases](https://github.com/m1101e/ansistated/releases/latest) page, extract it, and place the executable into your system `$PATH`:

```bash
# Example for Linux x86_64:
curl -sL https://github.com/m1101e/ansistated/releases/latest/download/ansistated_linux_amd64.tar.gz | tar -xz
sudo mv ansistated /usr/local/bin/
```

### Option 2. Via `go install`

If Go (1.21+) is installed on your system:

```bash
go install github.com/m1101e/ansistated@latest
```

### Option 3. Build from Source

```bash
git clone https://github.com/m1101e/ansistated.git
cd ansistated
go build -ldflags="-s -w" -o ansistated main.go
sudo mv ansistated /usr/local/bin/
```

---

## Configuration

By default, `ansistated` resolves the following paths:
* **Variables directory:** `~/.config/ansistated/vars/`
* **Templates directory:** `~/.config/ansistated/templates/`
* **Configuration file (optional):** `~/.config/ansistated/ansistated.yaml`

Example overrides via `ansistated.yaml`:
```yaml
vars_dir: "~/.config/ansistated/vars"
templates_dir: "~/.config/ansistated/templates"
```

Paths can also be provided directly via the `--vars-dir` and `--templates-dir` CLI flags.

---

## File Specifications

### 1. Variables (`vars/*.yaml`)

Variables are stored in standard YAML files. Files may include an optional frontmatter header assigning them to groups:

```yaml
---
group: theme, terminal
---
colors:
  background: "#1e1e2e"
  foreground: "#cdd6f4"
  accent: "rgb(137, 180, 250)"
```

### 2. Templates (`templates/*`)

Each template must contain frontmatter metadata followed by the template body:

```yaml
---
dest: "~/.config/alacritty/alacritty.toml"
group: terminal
post_cmd: "killall -USR1 alacritty || true"
---
[colors.primary]
background = "{{ colors.background | to_hex }}"
foreground = "{{ colors.foreground | to_hex }}"

[colors.cursor]
cursor = "{{ colors.accent | to_hex }}"
```

---

## CLI Usage

### Rendering (`render`)

Renders templates into their destination files. If the destination file already exists and its content matches the rendered output, the disk write is skipped.

```bash
# Render all templates
ansistated render

# Render templates belonging to a specific group
ansistated render --groups terminal

# Render multiple groups
ansistated render --groups terminal,wm
```

### Updating Variables (`update`)

Reads JSON input from `stdin` and updates values inside the YAML store on disk without stripping comments or altering node structure:

```bash
# Update targeted variables in a specific file (relative to vars/)
echo '{"colors": {"background": "#000000"}}' | ansistated update -f colors.yaml

# Update targeted variables across all files in a specific group
echo '{"colors": {"accent": "#fab387"}}' | ansistated update -g theme
```

---

## Color Filters

The following custom filters are available in templates for color format transformation:

| Filter | Description | Input Example | Output Example |
| :--- | :--- | :--- | :--- |
| `to_hex` | Normalizes color to `#RRGGBB` format | `rgb(255, 0, 0)` | `#ff0000` |
| `to_strip_hex` | Strips the leading `#` symbol | `#ff0000` | `ff0000` |
| `to_rgb` | Extracts `R, G, B` numeric triplet string | `#ffffff` | `255, 255, 255` |

---

## License

MIT
