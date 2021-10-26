# cms-installer (cmsi)
This is a command line tool & Go library to install CMS such as WordPress on a hosting server.
(CMS other than WordPress will be supported in the future)

## Example

```
CMSI_FTP_PASSOWRD="xxxxx" CMSI_DB_PASSOWRD="xxxxx" CMSI_SITE_PASSWORD="xxxxx" cmsi wp '{
"ftp_login_id": "your-ftp-id",
"ftp_host": "ftp.example.com",
"ftp_port": "21",
"ftp_dir": "www.example.com/blog",
"db_name": "your-db-name",
"db_user": "your-db-user",
"db_host": "your-db-host",
"db_prefix": "wp_",
"site_url": "https://www.example.com/blog",
"site_title": "MyDiary",
"site_user": "your-wp-user",
"site_email": "your@example.com"
}'

# Example of execution result
# 2021-10-26T03:56:49+09:00 INFO [19e1f3ab] Running https://www.example.com/blog - InjectInitScript
# 2021-10-26T03:56:50+09:00 INFO [19e1f3ab] Finished in 0.864 seconds (successful).
# 2021-10-26T03:56:50+09:00 INFO [4677c8a8] Running https://www.example.com/blog - ExecInit
# 2021-10-26T03:56:52+09:00 INFO [4677c8a8] Finished in 1.894 seconds (successful).
# 2021-10-26T03:56:52+09:00 INFO [94574975] Running https://www.example.com/blog - DeleteInitScript
# 2021-10-26T03:56:52+09:00 INFO [94574975] Finished in 0.384 seconds (successful).
# 2021-10-26T03:56:52+09:00 INFO [9158605a] Running https://www.example.com/blog - DeleteWpConfig
# 2021-10-26T03:56:53+09:00 INFO [9158605a] Finished in 0.505 seconds (successful).
# 2021-10-26T03:56:53+09:00 INFO [3dcd98d1] Running https://www.example.com/blog - WpAdminSetupConfig
# 2021-10-26T03:56:53+09:00 INFO [3dcd98d1] Finished in 0.24800001 seconds (successful).
# 2021-10-26T03:56:53+09:00 INFO [68cb313c] Running https://www.example.com/blog - WpAdminInstall
# 2021-10-26T03:56:55+09:00 INFO [68cb313c] Finished in 2.016 seconds (successful).
# Installation completed
```

## How it work
- Generate a PHP file (init.php) to download and extract the CMS archive, and upload it via FTP.
- Access init.php via http(s) to download and extract the CMS archive.
- POST to the CMS setup URL to complete the installation.
