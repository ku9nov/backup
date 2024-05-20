# Systemd configuration

Copy a service unit file. Service units installed by the system administrator are typically stored in `/etc/systemd/system/` directory, but this may vary depending on the Linux distribution.

The above service executes the `backup` command with `loglevel` and `config` environment variables.

Timer unit files contain information about a timer controlled and supervised by systemd. By default, a service with the same name as the timer is activated.

Copy a timer unit file in the same directory as the service file. The configuration below will activate the service everyday, at 5:30 AM.

## Using systemctl and journalctl
To start the timer:

```
sudo systemctl start backup.timer
```

To enable the timer to be started on boot-up:

```
sudo systemctl enable backup.timer
```

To show status information for the timer:

```
sudo systemctl status backup.timer
```

To show journal entries for the timer:

```
sudo journalctl -u backup.service
```