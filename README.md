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

Above, there is an example for my project, we use this logic:

On restic backup server we backup many environments and we have a rule for create backup repository: 

1. An environment is identified by `projectID` and `envId`:
    - `projectId`: is the directory name which contains `envId`, example `my-customer`
    - `envId`: is the directory of this environments, example `production`
    
2. Backup server set all `projectId` and `envId` in the same parents backup directory, example:
    - `/home/backup/vault/`in this directory we have all `projectId` and `envId` directory
    - in each `projectId/envId` we have a `.restic_confg_file` secret
    
3. Restic-exporter loop over each `projectId/envId` and get `.restic_config_file`, capture environment values and check restic repository to have metrics
    
```shell
/home/backup/vault/
|____my-projectId-1
| |____my-envId
| | |____.restic_config_file
|____my-projectId-2
| |____test
|____my-projectId-3
| |____dev
| | |____.restic_config_file

```

All files are used by restic-exporter on loop over. restic-exporter loop on .restic_config_file and get `repositoryUrl` and `repositoryPass` with regexp.

Each `.restic_config_file` must be have:

```shell
export RESTIC_REPOSITORY=resticRepoUrl
export RESTIC_PASSWORD=resticPassword
```

So define this environments values:

```shell
RESTIC_CREDENTIALS_PATH: /home/backup/vault/
RESTIC_CREDENTIALS_FILE: .restic_config_file
```

Now, you can let `restic-exporter` generate a `.prom` file which in turn is picked
up and exposed to prometheus by `node-exporter`.

This will generate many prometheus metrics:
```
restic_snapshot_timestamp{repository="projectID+envId", s3Bucket="s3://repoUrl"} 1.599849001e+09
restic_total_size{repository="projectID+envId", s3Bucket="s3://repoUrl"} 1.599849001e+09
restic_total_file_count{repository="projectID+envId", s3Bucket="s3://repoUrl"} 7.466395851e+09
restic_snapshot_count{repository="projectID+envId", s3Bucket="s3://repoUrl"} 2
```

## Why not make this a 'real' exporter

I went the `node-exporter` route here because I personally store backups on a NAS
with spinning hard drives that are spun down most of the time. Thus, I needed a way
to check the backup status at a specific time where the drives are already spun up.

The easiest way to achieve this was the file-generation route in combination with
a cronjob.
