package processor

import (
	"bytes"
	"os/exec"

	"github.com/sirupsen/logrus"
)

// Languages scanner支持的语言列表
var Languages []string = []string{
	"cs",
	"flex",
	"go",
	"java",
	"js",
	"php",
	"py",
	"ts",
	"web",
	"xml",
}

// Processor scanner processor
type Processor struct {
	Branch               string
	ProjectURL           string
	ProjectName          string
	ProjectKey           string
	Sources              string
	JSONDir              string
	JarDir               string
	SCMDisabled          string
	Language             string
	Rules                string
	KeepReport           string
	LogLevel             string
	Command              string
	Inclusions           string
	Exclusions           string
	GlobalExclusions     string
	GlobalTestExclusions string
	TestInclusions       string
	TestExclusions       string
	Logger               *logrus.Logger
}

// Scanner processor scanner
func (p *Processor) Scanner() {
	err := p.gitClone(p.ProjectURL, p.Branch, p.ProjectName)
	if err != nil {
		p.Logger.Errorln(err)
		return
	}
	args := make(map[string]string)
	args["sonar.projectKey"] = p.ProjectKey
	args["sonar.sources"] = p.Sources
	args["sonar.projectBaseDir"] = p.ProjectName
	args["sonar.jarDir"] = p.JarDir
	args["sonar.jsonDir"] = p.JSONDir
	args["sonar.ci.language"] = p.Language
	args["sonar.ci.rules"] = p.Rules
	args["sonar.scanner.keepReport"] = p.KeepReport
	args["sonar.scm.disabled"] = p.SCMDisabled
	args["sonar.log.level"] = p.LogLevel

	if p.Inclusions != "" {
		args["sonar.inclusions"] = p.Inclusions
	}
	if p.Exclusions != "" {
		args["sonar.exclusions"] = p.Exclusions
	}
	if p.GlobalExclusions != "" {
		args["sonar.global.exclusions"] = p.GlobalExclusions
	}
	if p.GlobalTestExclusions != "" {
		args["sonar.global.test.exclusions"] = p.GlobalTestExclusions
	}
	if p.TestInclusions != "" {
		args["sonar.test.inclusions"] = p.TestInclusions
	}
	if p.TestExclusions != "" {
		args["sonar.test.exclusions"] = p.TestExclusions
	}

	var params []string
	for k, v := range args {
		params = append(params, "-D"+k+"="+v)
	}
	cmd := exec.Command(p.Command, params...)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errLog bytes.Buffer
	cmd.Stderr = &errLog
	err = cmd.Run()
	if err != nil {
		p.Logger.Errorln(errLog.String())
		p.Logger.Errorln(err)
		return
	}
	p.Logger.Infoln(out.String())
}

func (p *Processor) gitClone(projectURL, branch, projectName string) error {
	cmd := exec.Command("git", "clone", "-b", branch, "--depth=1", projectURL, projectName)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errLog bytes.Buffer
	cmd.Stderr = &errLog
	err := cmd.Run()
	if err != nil {
		p.Logger.Errorln(errLog.String())
		return err
	}
	p.Logger.Infoln(out.String())
	p.Logger.Errorln(errLog.String())
	return nil
}
