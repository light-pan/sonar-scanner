package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

func setConfig(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	con.WorkerID = r.PostFormValue("worker_id")
	con.SiteID = r.PostFormValue("site_id")
	con.SiteLocation = r.PostFormValue("site_location")
	con.SiteIsp = r.PostFormValue("site_isp")
	con.Elasticsearch = r.PostFormValue("elasticsearch")
	res := result{
		State:   1,
		Message: "success",
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(res)
	w.Write(bytes)
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(con)
	w.Write(bytes)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	res := result{
		State:   1,
		Message: "success",
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(res)
	w.Write(bytes)
}

func scanner(w http.ResponseWriter, request *http.Request) {
	res := result{
		State:   1,
		Message: "success",
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	request.ParseMultipartForm(1024)
	data := request.PostFormValue("ext_data")
	extData := extData{}
	err := json.Unmarshal([]byte(data), &extData)
	if err != nil {
		res.State = 0
		res.Message = "Decode ext_data err."
		bytes, _ := json.Marshal(res)
		w.Write(bytes)
		return
	}

	p := &processor{
		Branch:               extData.Branch,
		ProjectURL:           extData.ProjectURL,
		ProjectName:          extData.ProjectName,
		ProjectKey:           extData.ProjectKey,
		Sources:              extData.Sources,
		JSONDir:              "/data/sonar-scanner/json",
		JarDir:               "/data/sonar-scanner/jar",
		Language:             extData.Language,
		Rules:                extData.Rules,
		Inclusions:           extData.Inclusions,
		Exclusions:           extData.Exclusions,
		GlobalExclusions:     extData.GlobalExclusions,
		GlobalTestExclusions: extData.GlobalTestExclusions,
		TestInclusions:       extData.TestInclusions,
		TestExclusions:       extData.TestExclusions,
		SCMDisabled:          "true",
		KeepReport:           "true",
		LogLevel:             "trace",
		Command:              "/data/sonar-scanner/bin/sonar-scanner",
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
			if !inArray(languages, language) {
				res.State = 0
				res.Message = "The languages should be in " + strings.Join(languages, ",")
				bytes, _ := json.Marshal(res)
				w.Write(bytes)
				return
			}
		}
	}
	if strings.TrimSpace(p.ProjectURL) == "" {
		res.State = 0
		res.Message = "The project_url must be not null"
		bytes, _ := json.Marshal(res)
		w.Write(bytes)
		return
	}
	if strings.TrimSpace(p.ProjectName) == "" {
		res.State = 0
		res.Message = "The project_name must be not null"
		bytes, _ := json.Marshal(res)
		w.Write(bytes)
		return
	}
	if strings.TrimSpace(p.ProjectKey) == "" {
		p.ProjectKey = p.ProjectName
	}

	go p.scanner()
	bytes, _ := json.Marshal(res)
	w.Write(bytes)
}

func inArray(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
