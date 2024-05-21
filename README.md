# backup

`Backup` is an application written in Golang designed to create backups intelligently and easily. This application supports multiple databases and multiple S3 storage providers. Additionally, it includes functionality to backup additional folders and files, ensuring proper configuration file preservation.

## Building the App:
You can build the app using the default Golang method

```
go build -o backup main.go
```

 or by running the `make` command.

`-loglevel` variable is used to configure log levels with possible values: debug, info, warn, error, fatal, panic.

`-config` sets the path to the config file with a default value of ./config.yml

## Compatibilities:
| Databases    | S3 Providers |
| -------- | ------- |
| MySQL :white_check_mark: | AWS S3   :white_check_mark: |
| PostgreSQL :white_check_mark: | Minio :white_check_mark: |
| MongoDB :white_check_mark: | Azure :white_check_mark: |
| Redis :x: | Digital Ocean :white_check_mark: |
| Elasticsearch :x: | GCP S3 :x: |

## Environment Overview:
You can review all possible variables in the [config.yml.example](config.yml.example) file.

## Systemd configuration
For more detailed information on systemd configuration and usage, please refer to the [README.md](examples/systemd/README.md) file.