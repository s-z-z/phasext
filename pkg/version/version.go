package version

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

var Ver *Version

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

func (v *Version) FullString() string {
	return fmt.Sprintf("%s.%s.%s", v.Major, v.Minor, v.Patch)
}

func (v *Version) Gte(major, minor, patch string) bool {
	return v.Major >= major && v.Minor >= minor && v.Patch >= patch
}

func (v *Version) Equal(major, minor, patch string) bool {
	return v.Major == major && v.Minor == minor && v.Patch == patch
}

func (v *Version) Lte(major, minor, patch string) bool {
	return v.Major <= major && v.Minor <= minor && v.Patch <= patch
}

func NewCmdVersion(v *Version) *cobra.Command {
	v.GitVersion = fmt.Sprintf("%s-%s.%s.%s", v.Branch, v.Major, v.Minor, v.Patch)
	v.GoVersion = runtime.Version()
	v.Compiler = runtime.Compiler
	v.Platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	Ver = v

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunVersion(os.Stdout, cmd, v)
		},
	}
	cmd.Flags().StringP("output", "o", "", "Output format; available util are 'yaml', 'json' and 'short'")
	return cmd
}

func RunVersion(out io.Writer, cmd *cobra.Command, v *Version) error {
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
