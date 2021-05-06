package main

import (
	"flag"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	output       = flag.String("output", "stats.txt", "File to export the stats to")
	resticBinary = flag.String("restic-bin", "restic", "Location of the restic binary to use (defaults to loading the one in your PATH)")
)

func collectMetrics(repoListConfig []configRepoRestic) *prometheus.Registry {
	registry := prometheus.NewRegistry()
	resticJobLabel := "restic-exporter"
	snapshot := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "restic_snapshot_timestamp",
	}, []string{"repository", "s3Bucket", "job"})

	err := registry.Register(snapshot)
	if err != nil {
		log.Fatal(err)
	}

	snapshotTotalSize := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "restic_snapshot_total_size",
	}, []string{"repository", "s3Bucket", "job"})
	err = registry.Register(snapshotTotalSize)
	if err != nil {
		log.Fatal(err)
	}
	snapshotTotalFile := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "restic_snapshot_total_file_count",
	}, []string{"repository", "s3Bucket", "job"})
	err = registry.Register(snapshotTotalFile)
	if err != nil {
		log.Fatal(err)
	}
	snapshotCount := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "restic_snapshot_count",
	}, []string{"repository", "s3Bucket", "job"})

	err = registry.Register(snapshotCount)
	if err != nil {
		log.Fatal(err)
	}

	for index, configItem := range repoListConfig {
		restic := Restic{Binary: *resticBinary, Name: configItem.RepositoryProjectId, Repository: configItem.RepositoryConf.RepositoryUrl, Password: configItem.RepositoryConf.RepositoryPass}
		timestamp, err := restic.SnapshotTimestamp()
		if err != nil {
			log.Printf("[%s] <ERR> %s", index, err)
		}
		totalSize, totalFileCount, err := restic.SnapshotsStats()
		if err != nil {
			log.Printf("[%s] <ERR> %s", index, err)
		}
		snapshotCounter, err := restic.SnapshotCount()
		if err != nil {
			log.Printf("[%s] <ERR> %s", index, err)
		}
		repository := configItem.RepositoryProjectId + "-" + configItem.RepositoryEnvId
		snapshot.WithLabelValues(repository, configItem.RepositoryConf.RepositoryUrl, resticJobLabel).Set(float64(timestamp))
		snapshotTotalSize.WithLabelValues(repository, configItem.RepositoryConf.RepositoryUrl, resticJobLabel).Set(float64(totalSize))
		snapshotTotalFile.WithLabelValues(repository, configItem.RepositoryConf.RepositoryUrl, resticJobLabel).Set(float64(totalFileCount))
		snapshotCount.WithLabelValues(repository, configItem.RepositoryConf.RepositoryUrl, resticJobLabel).Set(float64(snapshotCounter))
	}
	return registry
}

func main() {
	flag.Parse()

	getEnvVars := &EnvConfig{}
	envVars := getEnvVars.getEnvVars()
	repoListConfig := initResticConfigInMemory(envVars)
	registry := collectMetrics(repoListConfig)
	err := prometheus.WriteToTextfile(*output, registry)
	if err != nil {
		log.Fatal(err)
	}
}
