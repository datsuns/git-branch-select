// usage:
//	1.install to a searchable path
//  2. $ git branch-select
//
package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var verboseMode bool

var OptionFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "all, a",
		Usage: "select from --all",
	},
	cli.BoolFlag{
		Name:  "verbose, V",
		Usage: "verbose mode",
	},
}

func execute(tool string, params []string, debug bool) string {
	if debug {
		fmt.Printf(" >> %v %v\n", tool, params)
	}
	log, err := exec.Command(tool, params...).CombinedOutput()
	if err != nil {
		fmt.Println(string(log))
		panic(err)
	}
	return string(log)
}

func executable(bin string) bool {
	if path, err := exec.LookPath(bin); err != nil {
		fmt.Printf("%v\n", err)
		return false
	} else {
		if verboseMode {
			fmt.Printf("[%s] is located at [%s]\n", bin, path)
		}
		return true
	}
}

func generate_branch_list(from_all bool) []string {
	ret := []string{}
	var params [][]string
	if from_all {
		params = [][]string{{"branch", "--all"}}
	} else {
		params = [][]string{{"branch"}}
	}
	if !executable("git") {
		return ret
	}
	for _, p := range params {
		log := execute("git", p, verboseMode)
		for _, s := range strings.Split(log, "\n") {
			name := strings.Trim(s, " *")
			if len(name) > 0 {
				ret = append(ret, name)
			}
		}
	}
	return ret
}

// append "(?i)" means case insensitive search
func generate_branch_list_with_filter(from_all bool, filter string) []string {
	all := generate_branch_list(from_all)
	pattern := "(?i)" + fmt.Sprintf(".*%s.*", filter)
	key, _ := regexp.Compile(pattern)
	ret := []string{}
	for _, branch := range all {
		if key.FindStringIndex(branch) != nil {
			ret = append(ret, branch)
		}
	}

	return ret
}

func get_target_branch_index() (int, error) {
	fmt.Printf("insert target index: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(strings.Trim(text, "\r\n"))
}

func switch_git_branch(branch string) {
	params := [][]string{
		{"checkout", branch},
		{"status"},
	}
	for _, p := range params {
		log := execute("git", p, verboseMode)
		fmt.Println(log)
	}
}

func entry(c *cli.Context) error {
	var list []string
	var from_all = c.Bool("all")
	verboseMode = c.Bool("verbose")
	if c.NArg() == 0 {
		list = generate_branch_list(from_all)
	} else {
		fmt.Printf(" w/ filter [%s]\n", os.Args[1])
		list = generate_branch_list_with_filter(from_all, os.Args[1])
	}

	for i, b := range list {
		fmt.Printf("%2d: %s\n", i, b)
	}
	if num, err := get_target_branch_index(); err != nil {
		panic(err)
	} else {
		fmt.Printf("index is %d\n", num)
		switch_git_branch(list[num])
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "git-branch-select"
	app.Usage = "branch selector to checkout specified branch"
	app.Version = "1.0.0"
	app.Commands = nil
	app.Action = entry
	app.Flags = OptionFlags

	app.Run(os.Args)
}
