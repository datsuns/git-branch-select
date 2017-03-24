// usage:
//	1.install to a searchable path
//  2. $ git branch-select
//
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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
		fmt.Printf("[%s] is located at [%s]\n", bin, path)
		return true
	}
}

func generate_branch_list() []string {
	ret := []string{}
	params := [][]string{
		//{"branch", "--all"},
		{"branch"},
	}
	if !executable("git") {
		return ret
	}
	for _, p := range params {
		log := execute("git", p, true)
		for _, s := range strings.Split(log, "\n") {
			name := strings.Trim(s, " *")
			if len(name) > 0 {
				ret = append(ret, name)
			}
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
		log := execute("git", p, true)
		fmt.Println(log)
	}
}

func main() {
	list := generate_branch_list()

	for i, b := range list {
		fmt.Printf("%2d: %s\n", i, b)
	}
	if num, err := get_target_branch_index(); err != nil {
		panic(err)
	} else {
		fmt.Printf("index is %d\n", num)
		switch_git_branch(list[num])
	}
}
