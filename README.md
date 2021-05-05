# Restic exporter

Small prometheus exporter that does a single thing: It exports the time at which
the last [restic](github.com/restic/restic) backup has been done. That's it.

## Usage

This is not a standalone exporter - rather it generates `.prom` files that can
be picked up by [`node-exporter`](github.com/prometheus/node_exporter).

In order to use it, first setup `node-exporter` to watch a directory. If you are
on debian based systems, this can usually be done by adding the
`-collector.textfile.directory` option to `/etc/default/prometheus-node-exporter`:

```
ARGS="-collector.textfile.directory=\"/var/lib/prometheus/node-exporter/\""
```

## Restic configuration

For use restic-exporter you must be set  directory tree structure for get credentials for all restic repository.

Above, there is an example for my project

```shell
restic
|____resticRepo1
| |____.restic_config_file
|____resticRepo2
| |____.restic_config_file
|____resticRepo3
| |____.restic_config_file
|____resticRepo4
| |____.restic_config_file
```
All files are used by restic-exporter on loop over. restic-exporter loop on .restic_config_file and get `repositoryUrl` and `repositoryPass`

Above example for `.restic_config_file`, restic-exporter use regexp for capture only restic repository and password and use them on memory only, then for each restic repository we export all metrics.:

```shell
export RESTIC_REPOSITORY=resticRepoUrl
export RESTIC_PASSWORD=resticPassword
```

So define this environments values:

```shell
RESTIC_CREDENTIALS_PATH: /tmp/restic
RESTIC_CREDENTIALS_FILE: .restic_config_file
```

Now, you can let `restic-exporter` generate a `.prom` file which in turn is picked
up and exposed to prometheus by `node-exporter`.

This will generate many prometheus metrics:
```
restic_snapshot_timestamp{projectId="arbitrary-name-here", repository="s3://repoUrl"} 1.599849001e+09
restic_total_size{projectId="arbitrary-name-here", repository="s3://repoUrl"} 1.599849001e+09
restic_total_file_count{projectId="arbitrary-name-here", repository="s3://repoUrl"} 7.466395851e+09
```

## Why not make this a 'real' exporter

I went the `node-exporter` route here because I personally store backups on a NAS
with spinning hard drives that are spun down most of the time. Thus, I needed a way
to check the backup status at a specific time where the drives are already spun up.

The easiest way to achieve this was the file-generation route in combination with
a cronjob.
