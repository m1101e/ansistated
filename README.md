<div align="right">
  <a href="README.ru.md">Русский</a> | <b>English</b>
</div>

# ansistated

**ansistated** — a lightweight CLI templating engine and variable manager for dotfiles and configuration files.

> Originally developed to configure the [ansistate](https://github.com/m1101e/ansistate) desktop environment, it is designed with flexibility in mind and can be used with any templates.

---

## Features

- **Pongo2 Templating (Django/Jinja2 syntax)** with built-in filters for color format conversion (`to_hex`, `to_rgb`, `to_strip_hex`).


- **YAML Frontmatter in templates**: each template file defines its target path (`dest`), group (`group`), and post-processing command (`post_cmd`).


- **In-place YAML Updates (AST-based)**: the `update` command modifies variables from a JSON stream while preserving comments, indentation, and structure of the source YAML files.


- **Multithreaded Rendering**: parallel processing of templates with checksum verification to prevent unnecessary disk rewrites.


> **Merge Limitation (Key Conflicts):** You should not specify different scalar values for the same key across different files. Due to parallel file reading, key merging results may be unpredictable.


---

## Installation

### Option 1. Prebuilt Binaries

Download the archive for your platform from the [Releases](https://github.com/m1101e/ansistated/releases/latest) page, extract it, and move the binary to your executable path:

```bash
# Example for Linux x86_64:
curl -sL [https://github.com/m1101e/ansistated/releases/latest/download/ansistated_linux_amd64.tar.gz](https://github.com/m1101e/ansistated/releases/latest/download/ansistated_linux_amd64.tar.gz) | tar -xz
sudo mv ansistated /usr/local/bin/

```

### Option 2. Via `go install`

If Go is installed on your system:

```bash
go install [github.com/m1101e/ansistated@latest](https://github.com/m1101e/ansistated@latest)

```

### Option 3. Build from Source

```bash
git clone [https://github.com/m1101e/ansistated.git](https://github.com/m1101e/ansistated.git)
cd ansistated
go build -ldflags="-s -w" -o ansistated main.go
sudo mv ansistated /usr/local/bin/

```

---

## Configuration

By default, `ansistated` uses the following paths:

* **Variables directory:** `~/.config/ansistated/vars/`
* **Templates directory:** `~/.config/ansistated/templates/`
* **Config file (optional):** `~/.config/ansistated/ansistated.yaml`

Example override via `ansistated.yaml`:

```yaml
vars_dir: "~/.config/ansistated/vars"
templates_dir: "~/.config/ansistated/templates"

```

Paths can also be set via the `--vars-dir` and `--templates-dir` flags.

---

## File Format

### 1. Variables (`vars/*.yaml`)

Variables are stored in standard YAML files. They can include a header with groups:

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

Each template contains a frontmatter header with metadata followed by the template body:

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

### Reading Variables (`read`)

The `read` command outputs variable values to `stdout` in the selected format. Data is gathered from a specific file or merged from multiple files.

```bash
# Output all variables (from all files) in JSON format by default
ansistated read

# Output variables for a specific group in YAML format
ansistated read --group terminal --yaml

# Output variables from a specific file (relative to vars directory)
ansistated read --file theme/colors.yaml

```

### Rendering (`render`)

Generates files from templates. If the destination file already exists and hasn't changed, writing is skipped.

```bash
# Render all templates
ansistated render

# Render templates of a specific group
ansistated render --groups terminal

# Render multiple groups
ansistated render --groups terminal,wm

```

### Updating Variables (`update`)

Accepts JSON from `stdin` and updates values in YAML files without removing existing comments or altering structure:

```bash
# Update by relative path to the file
echo '{"colors": {"background": "#000000"}}' | ansistated update -f colors.yaml

# Update all files of a specific group
echo '{"colors": {"accent": "#fab387"}}' | ansistated update -g theme

```

---

## Color Filters

Built-in filters are available in templates for converting color formats:

| Filter | Description | Example Input | Example Output |
| --- | --- | --- | --- |
| `to_hex` | Converts color to `#RRGGBB` format | `rgb(255, 0, 0)` | `#ff0000` |
| `to_strip_hex` | Converts to hex without the `#` symbol | `#ff0000` | `ff0000` |
| `to_rgb` | Extracts `R, G, B` component values | `#ffffff` | `255, 255, 255` |

---

## License

MIT
