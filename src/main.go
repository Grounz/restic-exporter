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

	registry.Register(snapshot)

	snapshotTotalSize := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "restic_snapshot_total_size",
	}, []string{"repository", "s3Bucket", "job"})
	registry.Register(snapshotTotalSize)

	snapshotTotalFile := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "restic_snapshot_total_file_count",
	}, []string{"repository", "s3Bucket", "job"})
	registry.Register(snapshotTotalFile)

	for index, configItem := range repoListConfig {
		restic := Restic{Binary: *resticBinary, Name: configItem.RepositoryProjectId, Repository: configItem.RepositoryConf.RepositoryUrl, Password: configItem.RepositoryConf.RepositoryPass}
		timestamp, err := restic.SnapshotTimestamp()
		totalSize, totalFileCount, err := restic.SnapshotsStats()
		if err != nil {
			log.Printf("[%s] <ERR> %s", index, err)
		}
		repository := configItem.RepositoryProjectId + "-" + configItem.RepositoryEnvId
		snapshot.WithLabelValues(repository, configItem.RepositoryConf.RepositoryUrl, resticJobLabel).Set(float64(timestamp))
		snapshotTotalSize.WithLabelValues(repository, configItem.RepositoryConf.RepositoryUrl, resticJobLabel).Set(float64(totalSize))
		snapshotTotalFile.WithLabelValues(repository, configItem.RepositoryConf.RepositoryUrl, resticJobLabel).Set(float64(totalFileCount))
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
