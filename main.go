package main

import (
	"errors"
	"gopkg.in/xmlpath.v2"
	"log"
	"os"
	"path"
)
import (
	"github.com/urfave/cli/v2"
)

func main() {
	dto := &reqDto{}
	cw, _ := os.Getwd()
	_ = (&cli.App{
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Aliases:     []string{"p"},
				Value:       cw,
				Destination: &dto.path,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "version",
				Action: version(dto),
				Usage:  "Get project version",
			},
		},
		DefaultCommand: "version",
	}).Run(os.Args)
}

func version(dto *reqDto) cli.ActionFunc {
	return func(context *cli.Context) error {
		f, err := os.Open(dto.path)
		if stat, err := f.Stat(); err == nil && stat.IsDir() {
			return findProjectFile(dto, context)
		}
		return err
	}
}

func mavenProjectVersion(f *os.File) error {
	versionPath := xmlpath.MustCompile("/project/version")
	root, err := xmlpath.Parse(f)
	ifErr(err, "Failed to parse")
	if version, ok := versionPath.String(root); ok {
		writeString(version)
	} else {
		return errors.New("unable to retrieve maven project version")
	}
	return nil
}

func writeString(value string) {
	_, _ = os.Stdout.WriteString(value)
}

func findProjectFile(dto *reqDto, context *cli.Context) error {
	supportedNames := []string{"package.json", "pom.xml"}
	for i, name := range supportedNames {
		f, err := os.Open(path.Join(dto.path, name))
		if err == nil {
			switch i {
			case 1:
				return mavenProjectVersion(f)
			}
		}
	}
	return errors.New("project not supported")
}

func ifErr(err error, desc string) {
	if err != nil {
		log.Fatalf("%s: %s", desc, err.Error())
	}
}

type reqDto struct {
	path string
}
