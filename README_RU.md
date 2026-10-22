<div align="right">
  <b>Русский</b> | <a href="README.md">English</a>
</div>

# ansistated

**ansistated** — легковесный CLI-шаблонизатор и менеджер переменных для dotfiles и конфигурационных файлов.

> Утилита изначально разрабатывалась для конфигурации рабочего окружения [ansistate](https://github.com/m1101e/ansistate), однако спроектирована с оглядкой на универсальность и может использоваться с любыми шаблонами.

---

## Возможности

- **Шаблонизация Pongo2 (Django/Jinja2-синтаксис)** с встроенными фильтрами для конвертации форматов цвета (`to_hex`, `to_rgb`, `to_strip_hex`).
- **YAML Frontmatter в шаблонах**: каждый файл шаблона сам определяет целевой путь (`dest`), группу (`group`) и команду пост-обработки (`post_cmd`).
- **Точечное обновление YAML (AST-based)**: команда `update` модифицирует переменные из JSON-потока, сохраняя комментарии, отступы и структуру исходных YAML-файлов.
- **Многопоточный рендеринг**: параллельная обработка шаблонов с проверкой контрольных сумм для предотвращения лишних перезаписей на диск.

---

## Установка

### Вариант 1. Готовые бинарные файлы (Рекомендуется)

Скачайте архив для вашей платформы со страницы [Releases](https://github.com/m1101e/ansistated/releases/latest), распакуйте и переместите исполняемый файл в системный путь:

```bash
# Пример для Linux x86_64:
curl -sL https://github.com/m1101e/ansistated/releases/latest/download/ansistated_linux_amd64.tar.gz | tar -xz
sudo mv ansistated /usr/local/bin/
```

### Вариант 2. Через `go install`

Если в системе установлен Go (1.21+):

```bash
go install github.com/m1101e/ansistated@latest
```

### Вариант 3. Сборка из исходников

```bash
git clone https://github.com/m1101e/ansistated.git
cd ansistated
go build -ldflags="-s -w" -o ansistated main.go
sudo mv ansistated /usr/local/bin/
```

---

## Конфигурация

По умолчанию `ansistated` использует следующие пути:
* **Директория переменных:** `~/.config/ansistated/vars/`
* **Директория шаблонов:** `~/.config/ansistated/templates/`
* **Файл настроек (опционально):** `~/.config/ansistated/ansistated.yaml`

Пример переопределения через `ansistated.yaml`:
```yaml
vars_dir: "~/.config/ansistated/vars"
templates_dir: "~/.config/ansistated/templates"
```

Пути также можно передавать через флаги `--vars-dir` и `--templates-dir`.

---

## Формат файлов

### 1. Переменные (`vars/*.yaml`)

Переменные хранятся в обычных YAML-файлах. Они могут содержать заголовок с группами:

```yaml
---
group: theme, terminal
---
colors:
  background: "#1e1e2e"
  foreground: "#cdd6f4"
  accent: "rgb(137, 180, 250)"
```

### 2. Шаблоны (`templates/*`)

Каждый шаблон содержит заголовок frontmatter с метаданными и тело шаблона:

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

## Использование CLI

### Рендеринг (`render`)

Генерирует файлы из шаблонов. Если файл назначения уже существует и не изменился, запись пропускается.

```bash
# Рендеринг всех шаблонов
ansistated render

# Рендеринг шаблонов определенной группы
ansistated render --groups terminal

# Рендеринг нескольких групп
ansistated render --groups terminal,wm
```

### Обновление переменных (`update`)

Принимает JSON из `stdin` и обновляет значения в YAML-файлах, не удаляя существующие комментарии и структуру:

```bash
# Обновление по относительному пути к файлу
echo '{"colors": {"background": "#000000"}}' | ansistated update -f colors.yaml

# Обновление всех файлов определенной группы
echo '{"colors": {"accent": "#fab387"}}' | ansistated update -g theme
```

---

## Фильтры для работы с цветом

В шаблонах доступны встроенные фильтры для преобразования форматов цвета:

| Фильтр | Назначение | Пример входа | Пример выхода |
| :--- | :--- | :--- | :--- |
| `to_hex` | Приводит цвет к формату `#RRGGBB` | `rgb(255, 0, 0)` | `#ff0000` |
| `to_strip_hex` | Приводит к hex без символа `#` | `#ff0000` | `ff0000` |
| `to_rgb` | Извлекает значения компонентов `R, G, B` | `#ffffff` | `255, 255, 255` |

---

## Лицензия
MIT
