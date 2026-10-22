package plugins

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
)

type ColorFormat int

const (
	FormatUnknown ColorFormat = iota
	FormatHex
	FormatRgb
	FormatArgb
)

func detectColorType(val string) ColorFormat {
	val = strings.TrimSpace(val)
	if val == "" {
		return FormatUnknown
	}

	val = strings.TrimPrefix(val, "rgba(")
	val = strings.TrimPrefix(val, "rgb(")
	val = strings.TrimSuffix(val, ")")
	val = strings.TrimSpace(val)

	if strings.Contains(val, ",") {
		parts := strings.Split(val, ",")
		if len(parts) == 3 {
			return FormatRgb
		} else if len(parts) == 4 {
			return FormatArgb
		}
		return FormatUnknown
	}

	clean := strings.TrimPrefix(val, "#")
	if len(clean) == 3 || len(clean) == 4 || len(clean) == 6 || len(clean) == 8 {
		_, err := strconv.ParseInt(clean, 16, 64)
		if err == nil {
			return FormatHex
		}
	}

	return FormatUnknown
}

func convertToRgbString(val string) string {
	format := detectColorType(val)
	clean := strings.TrimSpace(val)
	clean = strings.TrimPrefix(clean, "rgba(")
	clean = strings.TrimPrefix(clean, "rgb(")
	clean = strings.TrimSuffix(clean, ")")
	clean = strings.TrimSpace(clean)

	switch format {
	case FormatRgb:
		return clean

	case FormatArgb:
		parts := strings.Split(clean, ",")
		if len(parts) == 4 {
			r := strings.TrimSpace(parts[0])
			g := strings.TrimSpace(parts[1])
			b := strings.TrimSpace(parts[2])
			return fmt.Sprintf("%s, %s, %s", r, g, b)
		}

	case FormatHex:
		cleanHex := strings.TrimPrefix(clean, "#")

		if len(cleanHex) == 3 {
			cleanHex = string([]byte{cleanHex[0], cleanHex[0], cleanHex[1], cleanHex[1], cleanHex[2], cleanHex[2]})
		} else if len(cleanHex) == 4 {
			cleanHex = string([]byte{cleanHex[0], cleanHex[0], cleanHex[1], cleanHex[1], cleanHex[2], cleanHex[2]})
		}

		if len(cleanHex) >= 6 {
			r, _ := strconv.ParseInt(cleanHex[0:2], 16, 64)
			g, _ := strconv.ParseInt(cleanHex[2:4], 16, 64)
			b, _ := strconv.ParseInt(cleanHex[4:6], 16, 64)
			return fmt.Sprintf("%d, %d, %d", r, g, b)
		}
	}

	return val
}

func convertToHex(val string) string {
	format := detectColorType(val)
	clean := strings.TrimSpace(val)
	clean = strings.TrimPrefix(clean, "rgba(")
	clean = strings.TrimPrefix(clean, "rgb(")
	clean = strings.TrimSuffix(clean, ")")
	clean = strings.TrimSpace(clean)

	switch format {
	case FormatHex:
		cleanHex := strings.TrimPrefix(clean, "#")
		if len(cleanHex) == 3 {
			cleanHex = string([]byte{cleanHex[0], cleanHex[0], cleanHex[1], cleanHex[1], cleanHex[2], cleanHex[2]})
		} else if len(cleanHex) == 8 {
			cleanHex = cleanHex[0:6]
		}
		return "#" + cleanHex

	case FormatRgb, FormatArgb:
		parts := strings.Split(clean, ",")
		if len(parts) >= 3 {
			r, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
			g, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
			b, err3 := strconv.Atoi(strings.TrimSpace(parts[2]))
			if err1 != nil || err2 != nil || err3 != nil {
				return val
			}
			return fmt.Sprintf("#%02x%02x%02x", r, g, b)
		}
	}

	return val
}

func init() {
	pongo2.RegisterFilter("to_rgb", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		str := in.String()
		result := convertToRgbString(str)
		return pongo2.AsValue(result), nil
	})

	pongo2.RegisterFilter("to_hex", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		str := in.String()
		result := convertToHex(str)
		return pongo2.AsValue(result), nil
	})

	pongo2.RegisterFilter("to_strip_hex", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		str := in.String()
		hexVal := convertToHex(str)
		stripped := strings.TrimPrefix(hexVal, "#")
		return pongo2.AsValue(stripped), nil
	})
}
