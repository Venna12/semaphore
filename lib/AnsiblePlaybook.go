package lib

import (
	"bytes"
	"fmt"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/util"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type AnsiblePlaybook struct {
	TemplateID int
	Repository db.Repository
	Logger     Logger
}

func (p AnsiblePlaybook) makeCmd(command string, args []string) *exec.Cmd {
	cmd := exec.Command(command, args...) //nolint: gas
	cmd.Dir = p.GetFullPath()

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s", util.Config.TmpPath))
	cmd.Env = append(cmd.Env, fmt.Sprintf("PWD=%s", cmd.Dir))
	cmd.Env = append(cmd.Env, fmt.Sprintln("PYTHONUNBUFFERED=1"))

	return cmd
}

func (p AnsiblePlaybook) runCmd(command string, args []string) error {
	cmd := p.makeCmd(command, args)
	p.Logger.LogCmd(cmd)
	return cmd.Run()
}

func (p AnsiblePlaybook) GetHosts(args []string) (hosts []string, err error) {
	args = append(args, "--list-hosts")
	cmd := p.makeCmd("ansible-playbook", args)

	var errb bytes.Buffer
	cmd.Stderr = &errb

	out, err := cmd.Output()
	if err != nil {
		return
	}

	re := regexp.MustCompile(`(?m)^\\s{6}(.*)$`)
	matches := re.FindAllSubmatch(out, 20)
	hosts = make([]string, len(matches))
	for i := range matches {
		hosts[i] = string(matches[i][1])
	}

	return
}

func (p AnsiblePlaybook) RunPlaybook(args []string, cb func(*os.Process)) error {
	cmd := p.makeCmd("ansible-playbook", args)
	p.Logger.LogCmd(cmd)
	cmd.Stdin = strings.NewReader("")
	err := cmd.Start()
	if err != nil {
		return err
	}
	cb(cmd.Process)
	return cmd.Wait()
}

func (p AnsiblePlaybook) RunGalaxy(args []string) error {
	return p.runCmd("ansible-galaxy", args)
}

func (p AnsiblePlaybook) GetFullPath() (path string) {
	path = p.Repository.GetFullPath(p.TemplateID)
	return
}
