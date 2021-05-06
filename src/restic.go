package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type StatsResponse struct {
	TotalSize      int64 `json:"total_size"`
	TotalFileCount int64 `json:"total_file_count"`
}

type SnapshotResponse struct {
	Time string `json:"time"`
}

type SnapshotCount interface{}

type Restic struct {
	Binary     string
	Name       string
	Repository string
	Password   string
}

func (restic Restic) Run(arguments []string, target interface{}) error {
	arguments = append(arguments, "--json")

	log.Printf("[%s] %s %s", restic.Name, restic.Binary, arguments)
	command := exec.Command(restic.Binary, arguments...)
	command.Env = append(os.Environ(), fmt.Sprintf("RESTIC_REPOSITORY=%s", restic.Repository), fmt.Sprintf("RESTIC_PASSWORD=%s", restic.Password))
	output, err := command.Output()
	if err != nil {
		return err
	}
	err = json.Unmarshal(output, target)
	if err != nil {
		return err
	}

	return nil
}

func (restic Restic) SnapshotCount() (int64, error) {
	snapshotCounter := make([]SnapshotCount, 0)
	err := restic.Run([]string{"snapshots"}, &snapshotCounter)
	if err != nil {
		return -1, err
	}
	snapCounter := len(snapshotCounter)
	return int64(snapCounter), err
}

func (restic Restic) SnapshotTimestamp() (int64, error) {
	snapshots := make([]SnapshotResponse, 0)
	err := restic.Run([]string{"snapshots", "latest"}, &snapshots)
	if err != nil {
		return -1, err
	}

	snapTime, err := time.Parse(time.RFC3339Nano, snapshots[0].Time)
	return snapTime.Unix(), err
}

func (restic Restic) SnapshotsStats() (int64, int64, error) {
	statsResponse := StatsResponse{}
	err := restic.Run([]string{"stats", "latest"}, &statsResponse)
	if err != nil {
		return -1, -1, err
	}
	return statsResponse.TotalSize, statsResponse.TotalFileCount, err
}
