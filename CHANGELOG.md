## 0.5.3
### Bug Fixes
- Fixed an issue where the mysqldump command failed due to special characters in the password.

## 0.5.2
### Features
- Added `dumpFlags` field for PostgreSQL and MySQL configuration.
- Added `Redis` backups.

## 0.5.1
### Bug Fixes
- Disabled CGO for the binary file to function correctly after building with GitHub Actions.

## 0.5.0
### Features
- :tada: **Implemented basic functionality.**
  - The application now includes basic features ready for testing in real environments and servers.
  - Added support for MySQL, MongoDB, and PostgreSQL databases.
  - Added support for additional backup file formats.
  - Added support for AWS S3, Azure Blob Storage, Minio, and Digital Ocean Spaces.
  - Added functionality to remove old backups after a specified time interval for cost and memory efficiency.
  - Added systemd configurations for scheduled backups.
  - Added functionality to create extra backups that will be uploaded to a separate bucket of another cloud provider for more reliable backup storage in multiple locations.
  - Implemented behavior such that if any of the backups fail, old files are not deleted.
  - Added support for Slack notifications and Zabbix sender.

## 0.0.1 (Example)

### Important Notes

- **Important**: Removed something
- **Important**: Updated something
- App now requires something

### Maintenance

- Removed a redundant feature
- Added a new functionality
- Improved overall performance
    - Enhanced user interface responsiveness

### Features

- :tada: Implemented a new feature
    - This feature allows users to...
- Added a user profile customization option

### Bug Fixes

- Fixed a critical issue that caused...