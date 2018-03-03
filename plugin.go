package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const kubectl = "/usr/local/bin/kubectl"

type (
	Repo struct {
		Owner   string
		Name    string
		Link    string
		Avatar  string
		Branch  string
		Private bool
		Trusted bool
	}

	Build struct {
		Number   int
		Event    string
		Status   string
		Deploy   string
		Created  int64
		Started  int64
		Finished int64
		Link     string
	}

	Commit struct {
		Remote  string
		Sha     string
		Ref     string
		Tag     string
		Link    string
		Branch  string
		Message string
		Author  Author
	}

	Author struct {
		Name   string
		Email  string
		Avatar string
	}

	Config struct {
		Manifest  string
		Namespace string
		Server    string
		Token     string
		Selector  string
		Prune     bool
		Ca        string
		capath    string
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Commit Commit
		Config Config
		Env    map[string]string
		cmds   []*exec.Cmd
	}
)

func (p *Plugin) Exec() error {
	if p.Config.Ca != "" {
		ca, err := base64.StdEncoding.DecodeString(p.Config.Ca)
		if err != nil {
			return err
		}
		tmpfile, err := ioutil.TempFile("", "k8s_ca")
		if err != nil {
			return err
		}
		defer os.Remove(tmpfile.Name())
		tmpfile.Write(ca)
		tmpfile.Sync()
		p.Config.capath = tmpfile.Name()
	}

	fmt.Printf("+ drone-k8s - %s (%s) %s\n", version, commit, date)

	t, err := template.ParseFiles(p.Config.Manifest)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	p.Env = map[string]string{}
	for _, setting := range os.Environ() {
		pair := strings.SplitN(setting, "=", 2)
		p.Env[pair[0]] = pair[1]
	}
	if err := t.Execute(&buf, p); err != nil {
		return err
	}
	if err := p.commandVersion(); err != nil {
		return err
	}
	if err := p.commandApply(&buf); err != nil {
		return err
	}
	fmt.Printf("%v\n", p.cmds)
	for _, cmd := range p.cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		o := strings.Join(cmd.Args, " ")
		if p.Config.Token != "" {
			o = strings.Replace(o, p.Config.Token, "********", -1)
		}
		fmt.Printf("+ %s\n", o)

		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Plugin) commandVersion() error {
	p.cmds = append(p.cmds, exec.Command(kubectl, "version", "--client"))
	return nil
}

func (p *Plugin) commandApply(buf *bytes.Buffer) error {
	args := []string{"apply", "-f", "-"}
	if p.Config.Prune && p.Config.Selector != "" {
		args = append(args, "--prune", "-l", p.Config.Selector)
	}
	if p.Config.Namespace != "" {
		args = append(args, "--namespace="+p.Config.Namespace)
	}
	if p.Config.Server != "" {
		args = append(args, "--server="+p.Config.Server)
	}
	if p.Config.Token != "" {
		args = append(args, "--token="+p.Config.Token)
	}
	if p.Config.Ca != "" {
		args = append(args, "--certificate-authority="+p.Config.capath)
	}
	cmd := exec.Command(kubectl, args...)
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer in.Close()
		buf.WriteTo(in)
	}()
	p.cmds = append(p.cmds, cmd)
	return nil
}
