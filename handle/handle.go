package handle

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/light-pan/sonar-scanner/processor"
	"github.com/sirupsen/logrus"
)

type result struct {
	State   int    `json:"state"`
	Message string `json:"msg"`
}

type extData struct {
	Branch               string `json:"branch,omitempty"`
	ProjectURL           string `json:"project_url,omitempty"`
	ProjectName          string `json:"project_name,omitempty"`
	ProjectKey           string `json:"project_key,omitempty"`
	Sources              string `json:"sources,omitempty"`
	Language             string `json:"language,omitempty"`
	Inclusions           string `json:"inclusions,omitempty"`
	Exclusions           string `json:"exclusions,omitempty"`
	GlobalExclusions     string `json:"global_exclusions,omitempty"`
	GlobalTestExclusions string `json:"global_test_exclusions,omitempty"`
	TestInclusions       string `json:"test_inclusions,omitempty"`
	TestExclusions       string `json:"test_exclusions,omitempty"`
	Rules                string `json:"rules,omitempty"`
}

type config struct {
	WorkerID      string `json:"worker_id"`
	SiteID        string `json:"site_id"`
	SiteLocation  string `json:"site_location"`
	SiteIsp       string `json:"site_isp"`
	Elasticsearch string `json:"elasticsearch"`
}

// Handle http handle
type Handle struct {
	con      config
	Logger   *logrus.Logger
	JSONDir  string
	JarDir   string
	LogLevel string
	Command  string
}

// SetConfig 设置config
func (h *Handle) SetConfig(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	h.con.WorkerID = r.PostFormValue("worker_id")
	h.con.SiteID = r.PostFormValue("site_id")
	h.con.SiteLocation = r.PostFormValue("site_location")
	h.con.SiteIsp = r.PostFormValue("site_isp")
	h.con.Elasticsearch = r.PostFormValue("elasticsearch")
	res := result{
		State:   1,
		Message: "success",
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(res)
	w.Write(bytes)
}

// GetConfig 获取config
func (h *Handle) GetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	bytes, _ := json.Marshal(h.con)
	w.Write(bytes)
}

// DeleteTask 删除task
func (h *Handle) DeleteTask(w http.ResponseWriter, r *http.Request) {
	res := result{
		State:   1,
		Message: "success",
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(res)
	w.Write(bytes)
}

// Scanner 执行扫描任务
func (h *Handle) Scanner(w http.ResponseWriter, request *http.Request) {
	res := result{
		State:   1,
		Message: "success",
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	request.ParseMultipartForm(1024)
	data := request.PostFormValue("ext_data")
	p, err := transExtDataToProcess(data)
	if err != nil {
		res.State = 0
		res.Message = err.Error()
		bytes, _ := json.Marshal(res)
		w.Write(bytes)
		return
	}
	p.Logger = h.Logger
	p.JSONDir = h.JSONDir
	p.JarDir = h.JarDir
	p.Command = h.Command
	p.LogLevel = h.LogLevel

	go p.Scanner()
	bytes, _ := json.Marshal(res)
	w.Write(bytes)
}

func transExtDataToProcess(data string) (*processor.Processor, error) {
	eData := extData{}
	err := json.Unmarshal([]byte(data), &eData)
	if err != nil {
		return nil, errors.New("Decode ext_data err")
	}

	p := &processor.Processor{
		Branch:               eData.Branch,
		ProjectURL:           eData.ProjectURL,
		ProjectName:          eData.ProjectName,
		ProjectKey:           eData.ProjectKey,
		Sources:              eData.Sources,
		Language:             eData.Language,
		Rules:                eData.Rules,
		Inclusions:           eData.Inclusions,
		Exclusions:           eData.Exclusions,
		GlobalExclusions:     eData.GlobalExclusions,
		GlobalTestExclusions: eData.GlobalTestExclusions,
		TestInclusions:       eData.TestInclusions,
		TestExclusions:       eData.TestExclusions,
		SCMDisabled:          "true",
		KeepReport:           "true",
	}

	if strings.TrimSpace(p.Branch) == "" {
		p.Branch = "master"
	}
	if strings.TrimSpace(p.Sources) == "" {
		p.Sources = "."
	}
	lags := strings.Split(p.Language, ",")
	if strings.TrimSpace(p.Language) != "" {
		for _, language := range lags {
			if !inArray(processor.Languages, language) {
				return nil, errors.New("The languages should be in " + strings.Join(processor.Languages, ","))
			}
		}
	}
	if strings.TrimSpace(p.ProjectURL) == "" {
		return nil, errors.New("The project_url must be not null")
	}
	if strings.TrimSpace(p.ProjectName) == "" {
		return nil, errors.New("The project_name must be not null")
	}
	if strings.TrimSpace(p.ProjectKey) == "" {
		p.ProjectKey = p.ProjectName
	}

	return p, nil
}

func inArray(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
