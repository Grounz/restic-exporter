package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

type configRepoRestic struct {
	RepositoryName string
	RepositoryConf struct {
		RepositoryUrl string `yaml:"repository"`
		RepositoryPass string `yaml:"password"`
	}
}

type EnvConfig struct {
	resticCredentialsPathDirectory string
	resticCredentialsFile string
}

func (e EnvConfig) getEnvVars() EnvConfig {
	e.resticCredentialsPathDirectory = os.Getenv("RESTIC_CREDENTIALS_PATH")
	if e.resticCredentialsPathDirectory == "" {
		log.Print("No directory set for restic credentials, set to /root by default")
		e.resticCredentialsPathDirectory = "/root"
	}

	e.resticCredentialsFile = os.Getenv("RESTIC_CREDENTIALS_FILE")
	if e.resticCredentialsFile == "" {
		log.Print("No file name defined for restic credentials file, set to .restic_config_file by default")
		e.resticCredentialsFile = ".restic_config_file"
	}
	return e
}


func initResticConfigInMemory(envVars EnvConfig) []configRepoRestic{
	c := &configRepoRestic{}
	var repoList []configRepoRestic
	directoryList, err := ioutil.ReadDir(envVars.resticCredentialsPathDirectory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range directoryList {
		resticRepoName := file.Name()
		resticConfFile, err := os.Open(envVars.resticCredentialsPathDirectory + "/" + resticRepoName + "/" + envVars.resticCredentialsFile)
		if err != nil {
			log.Fatal(err)
		}
		defer resticConfFile.Close()
		scanner := bufio.NewScanner(resticConfFile)
		for scanner.Scan() {
			rS3RepoUrl := regexp.MustCompile(`^(?P<exportString>\w{6}) (?P<varRepo>RESTIC_REPOSITORY=)(?P<S3Url>s3:.+)`)
			rResticRepoPass := regexp.MustCompile(`^(?P<exportString>\w{6}) (?P<varPass>RESTIC_PASSWORD=)(?P<resticPass>.+)`)
			envRepoUrl := rS3RepoUrl.FindStringSubmatch(scanner.Text())
			envRepoPass := rResticRepoPass.FindStringSubmatch(scanner.Text())
			if len(envRepoUrl) >= 3 {
				c.RepositoryConf.RepositoryUrl = envRepoUrl[3]
			}
			if len(envRepoPass) >= 3 {
				c.RepositoryConf.RepositoryPass = envRepoPass[3]
			}
			c.RepositoryName = resticRepoName
		}
		repoList = append(repoList, *c)
	}
	return repoList
}
