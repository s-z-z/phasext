package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	gitBranch string
	verMajor  string
	verMinor  string
	verPatch  string
	gitURL    string
	gitCommit string
	buildDate string
)

type Version struct {
	GitVersion string `json:"gitVersion"`
	Branch     string `json:"branch"`
	Major      string `json:"major"`
	Minor      string `json:"minor"`
	Patch      string `json:"patch"`
	GitURL     string `json:"gitURL"`
	GitCommit  string `json:"gitCommit"`
	BuildDate  string `json:"buildDate"`
	GoVersion  string `json:"goVersion"`
	Compiler   string `json:"compiler"`
	Platform   string `json:"platform"`
}

func Get() Version {
	return Version{
		GitVersion: fmt.Sprintf("%s-v%s.%s.%s", gitBranch, verMajor, verMinor, verPatch),
		Branch:     gitBranch,
		Major:      verMajor,
		Minor:      verMinor,
		Patch:      verPatch,
		GitURL:     gitURL,
		GitCommit:  gitCommit,
		BuildDate:  buildDate,
		GoVersion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func newCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunVersion(os.Stdout, cmd)
		},
	}
	cmd.Flags().StringP("output", "o", "", "Output format; available util are 'yaml', 'json' and 'short'")
	return cmd
}

func RunVersion(out io.Writer, cmd *cobra.Command) error {
	v := Get()
	const flag = "output"
	of, err := cmd.Flags().GetString(flag)
	if err != nil {
		return errors.Wrapf(err, "error accessing flag %s for command %s", flag, cmd.Name())
	}

	var outStr string

	switch of {
	case "short":
		outStr = fmt.Sprintf("%s\n", v.GitVersion)
	case "json":
		y, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			return err
		}
		outStr = string(y)
	default:
		y, err := yaml.Marshal(&v)
		if err != nil {
			return err
		}
		outStr = string(y)
	}
	_, err = fmt.Fprintln(out, outStr)
	return err
}
