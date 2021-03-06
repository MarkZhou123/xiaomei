package misc

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	cmdPkg "github.com/lovego/cmd"
	"github.com/lovego/config/conf"
	"github.com/lovego/xiaomei/misc/dbs"
	"github.com/lovego/xiaomei/misc/godoc"
	"github.com/lovego/xiaomei/misc/token"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func Cmds(rootCmd *cobra.Command) []*cobra.Command {
	return append(
		dbs.Cmds(), godoc.Cmd(), timestampSignCmd(), token.Cmd(),
		specCmd(), coverCmd(), yamlCmd(), bashCompletionCmd(rootCmd),
	)
}

func timestampSignCmd() *cobra.Command {
	var secret string
	cmd := &cobra.Command{
		Use:   `timestamp-sign [<env>]`,
		Short: `Generate Timestamp and Sign headers for curl command.`,
		RunE: release.EnvCall(func(env string) error {
			ts := time.Now().Unix()
			if secret == "" {
				secret = release.AppConf(env).Secret
			}
			fmt.Printf("-H Timestamp:%d -H Sign:%s\n", ts, conf.TimestampSign(ts, secret))
			return nil
		}),
	}
	cmd.Flags().StringVarP(&secret, `secret`, `s`, ``, `secret used to generate sign`)
	return cmd
}

func coverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `cover [package] ...`,
		Short: `Show coverage details for packages.`,
		RunE: func(_ *cobra.Command, args []string) error {
			_, err := cmdPkg.Run(cmdPkg.O{}, "sh", "-c", fmt.Sprintf(`
rm -f /tmp/go-cover.out && {
  go test -p 1 --gcflags=-l -coverprofile /tmp/go-cover.out %s
  test -f /tmp/go-cover.out && (($(wc -c </tmp/go-cover.out) > 10)) && {
    go tool cover -func /tmp/go-cover.out | tail -n 1
    go tool cover -html /tmp/go-cover.out
  }
}`, strings.Join(args, " ")))
			return err
		},
	}
	return cmd
}

func yamlCmd() *cobra.Command {
	var goSyntax bool
	cmd := &cobra.Command{
		Use:   `yaml`,
		Short: `Parse yaml file.`,
		RunE: release.Arg1Call(``, func(p string) error {
			content, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}
			data := make(map[string]interface{})
			if err := yaml.Unmarshal(content, data); err != nil {
				return err
			}
			if goSyntax {
				fmt.Printf("%#v\n", data)
			} else {
				if buf, err := yaml.Marshal(data); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%s\n", buf)
				}
			}
			return nil
		}),
	}
	cmd.Flags().BoolVarP(&goSyntax, `go-syntax`, `g`, false, `print in go syntax`)
	return cmd
}
