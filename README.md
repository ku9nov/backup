# backup

`Backup` is an application written in Golang designed to create backups intelligently and easily. This application supports multiple databases and multiple S3 storage providers. Additionally, it includes functionality to backup additional folders and files, ensuring proper configuration file preservation.

## One command download 
```
BACKUP=$(curl --silent "https://api.github.com/repos/ku9nov/backup/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'); wget https://github.com/ku9nov/backup/releases/download/$BACKUP/backup -O /usr/bin/backup && chmod +x /usr/bin/backup
```

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
| Redis :white_check_mark: | Digital Ocean :white_check_mark: |
|  | GCP S3 :x: |

## Environment Overview:
You can review all possible variables in the [config.yml.example](config.yml.example) file.

## Systemd configuration
For more detailed information on systemd configuration and usage, please refer to the [README.md](examples/systemd/README.md) file.

## Enable Slack Notification
### Create Slack Bot
The first thing we need to do is to create the Slack application. Visit the [slack website](https://api.slack.com/apps?new_app=1) to create the application. Select the `From scratch` option. 
You will be presented with the option to add a Name to the application and the Workspace to allow the application to be used. You should be able to see all workspaces that you are connected to. Select the appropriate workspace.
Select the Bot option.
After clicking Bots you will be redirected to a Help information page, select the option to add scopes. The first thing we need to add to the application is the actual permissions to perform anything.
After pressing `Review Scopes to Add`, scroll down to Bot scopes and start adding the 4 scopes:
```
channels:history
chat:write
chat:write.customize
incoming-webhook
```
After adding the scopes we are ready to install the application. Once you click Allow you will see long strings, one OAuth token, and one Webhook URL. Remember the location of these, or save them on another safe storage. Then we need to invite the Application into a channel that we want him to be available in.
Go there and start typing a command message which is done by starting the message with `/`. We can invite the bot by typing `/invite @NameOfYourbot`.

## Local development
If you only want to run dependency services (minio), use this command:
```
docker-compose up
```

Then, open `http://localhost:9011/access-keys`, create Access Keys, and set `accessKey` and `secretKey` in either `config.yml`. To access the Minio dashboard, use the `MINIO_ROOT_USER` from the `docker-compose.yaml` as the username and `MINIO_ROOT_PASSWORD` as the password.