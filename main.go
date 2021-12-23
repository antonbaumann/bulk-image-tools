package main

import (
	"fmt"
	"github.com/antonbaumann/bulk-image-tools/internal"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func isValidRotation(r int) bool {
	return r == int(internal.Rotate0) ||
		r == int(internal.Rotate90) ||
		r == int(internal.Rotate180) ||
		r == int(internal.Rotate270)
}

func main() {
	app := &cli.App{
		Name:  "Bulk Image Tools",
		Usage: "Apply image transformation in bulk",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "from",
				Usage:    "source folder",
				Required: true,
			},
			&cli.PathFlag{
				Name:     "to",
				Usage:    "destination folder",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "width",
				Aliases: []string{"W"},
				Usage:   "Set width of all images to `width` preserving aspect ratio",
			},
			&cli.IntFlag{
				Name:    "height",
				Aliases: []string{"H"},
				Usage:   "Set height of all images to `height` preserving aspect ratio",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"F"},
				Usage:   "Convert transformed image to `format`",
			},
			&cli.IntFlag{
				Name:    "rotation",
				Aliases: []string{"R"},
				Usage:   "Rotate all images by `rotation`",
			},
			&cli.IntFlag{
				Name:        "workers",
				Usage:       "number of workers used",
				DefaultText: "4",
			},
		},
		Action: func(c *cli.Context) error {
			var format internal.ImageFormat
			if c.IsSet("format") {
				format = internal.ImageFormat(c.String("format"))
				if !format.IsValid() {
					return fmt.Errorf("%v is not a valid format", format)
				}
			} else {
				format = internal.JPEG
			}

			var rotation internal.ImageRotation
			if c.IsSet("rotation") {
				rotation = internal.ImageRotation(c.Int("rotation"))
				if !rotation.IsValid() {
					return fmt.Errorf("%vdeg is not a valid rotation", rotation)
				}
			} else {
				rotation = internal.Rotate0
			}

			workers := 4
			if c.IsSet("workers") {
				workers = c.Int("workers")
			}

			return internal.RunTransformations(
				c.Path("from"),
				c.Path("to"),
				format,
				c.Int("width"),
				c.Int("height"),
				rotation,
				workers,
			)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
