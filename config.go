package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"
)

// Config defines application settings.
type Config struct {
	Input      string  // Image file with simulation data to load.
	Width      int     // Display width in pixels.
	Height     int     // Display height in pixels.
	Palette    Palette // Color palette to use.
	Fullscreen bool    // Run in fullscreen mode?
}

// parseArgs parses commandline arguments and returns a config struct.
// Exits the program with an error if invalid data was found.
func parseArgs() *Config {
	var c Config
	c.Width = 1280
	c.Height = 600
	c.Fullscreen = false
	c.Palette.LoadDefault()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] <image>\n", os.Args[0])
		flag.PrintDefaults()
	}

	palEmpty := flag.String("pal-empty", hexStr(c.Palette.Empty), "Color for empty cells.")
	palWire := flag.String("pal-wire", hexStr(c.Palette.Wire), "Color for wire cells.")
	palHead := flag.String("pal-head", hexStr(c.Palette.Head), "Color for electron head cells.")
	palTail := flag.String("pal-tail", hexStr(c.Palette.Tail), "Color for electron tail cells.")

	flag.IntVar(&c.Width, "width", c.Width, "Display width in pixels.")
	flag.IntVar(&c.Height, "height", c.Height, "Display height in pixels.")
	flag.BoolVar(&c.Fullscreen, "fullscreen", c.Fullscreen, "Use a fullscreen display.")
	version := flag.Bool("version", false, "Displays version information.")
	flag.Parse()

	if *version {
		fmt.Println(Version())
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "missing input image")
		flag.Usage()
		os.Exit(1)
	}

	c.Input = flag.Arg(0)

	if c.Width <= 0 {
		fmt.Fprintf(os.Stderr, "width must be > 0")
		flag.Usage()
		os.Exit(1)
	}

	if c.Height <= 0 {
		fmt.Fprintf(os.Stderr, "height must be > 0")
		flag.Usage()
		os.Exit(1)
	}

	if len(*palEmpty) > 0 {
		c.Palette.Empty = parseHex(*palEmpty)
	}

	if len(*palWire) > 0 {
		c.Palette.Wire = parseHex(*palWire)
	}

	if len(*palHead) > 0 {
		c.Palette.Head = parseHex(*palHead)
	}

	if len(*palTail) > 0 {
		c.Palette.Tail = parseHex(*palTail)
	}

	return &c
}

// hexStr returns color c as a hex string.
// E.g.: [255, 255, 255] -> "ffffff"
// E.g.: [255, 0, 127] -> "ff007f"
func hexStr(clr color.RGBA) string {
	return fmt.Sprintf("%02x%02x%02x", clr.R, clr.G, clr.B)
}

// parseHex returns the given hex string as a floating-point RGBA color.
// E.g.: "ffffff" -> [255, 255, 255]
// E.g.: "ff007f" -> [255, 0, 127]
func parseHex(str string) color.RGBA {
	str = strings.ToLower(str)
	if len(str) != 6 {
		fmt.Fprintf(os.Stderr, "invalid color %q; expected form: rrggbb\n", str)
		os.Exit(1)
	}

	sr := str[0:2]
	sg := str[2:4]
	sb := str[4:6]

	r, err := strconv.ParseInt(sr, 16, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid red component in color %q; %v\n", str, err)
		os.Exit(1)
	}

	g, err := strconv.ParseInt(sg, 16, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid green component in color %q; %v\n", str, err)
		os.Exit(1)
	}

	b, err := strconv.ParseInt(sb, 16, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid blue component in color %q; %v\n", str, err)
		os.Exit(1)
	}

	return color.RGBA{byte(r), byte(g), byte(b), 255}
}
