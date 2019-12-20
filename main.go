package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type result struct {
	State   int    `json:"state"`
	Message string `json:"msg"`
}

func main() {
	http.HandleFunc("/scanner", scanner)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func scanner(w http.ResponseWriter, request *http.Request) {
	res := result{
		State:   1,
		Message: "success",
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if request.Method != http.MethodPost {
		res.State = 0
		res.Message = "The request method  must be post"
		bytes, _ := json.Marshal(res)
		w.Write(bytes)
		return
	}
	request.ParseMultipartForm(1024)
	branch := request.PostFormValue("branch")
	projectURL := request.PostFormValue("project_url")
	projectName := request.PostFormValue("project_name")
	projectKey := request.PostFormValue("project_key")
	sources := request.PostFormValue("sources")
	language := request.PostFormValue("language")
	inclusions := request.PostFormValue("inclusions")
	exclusions := request.PostFormValue("exclusions")
	globalExclusions := request.PostFormValue("global_exclusions")
	globalTestExclusions := request.PostFormValue("global_test_exclusions")
	testInclusions := request.PostFormValue("test_inclusions")
	testExclusions := request.PostFormValue("test_exclusions")
	rules := request.PostFormValue("rules")

	p := &processor{
		Branch:               branch,
		ProjectURL:           projectURL,
		ProjectName:          projectName,
		ProjectKey:           projectKey,
		Sources:              sources,
		JSONDir:              "/data/sonar-scanner/json",
		JarDir:               "/data/sonar-scanner/jar",
		Language:             language,
		Rules:                rules,
		Inclusions:           inclusions,
		Exclusions:           exclusions,
		GlobalExclusions:     globalExclusions,
		GlobalTestExclusions: globalTestExclusions,
		TestInclusions:       testInclusions,
		TestExclusions:       testExclusions,
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
