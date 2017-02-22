# Ftpbeat
Fully customizable Beat for FTP Server - this beat can ship the row of remote file through FTP or SFTP, or get files from remote server.


## Current status
Ftpbeat still on beta.

#### To Do:
* Support for ssh key 
* Support encrypted password


## Features
* Get files from remote server
* Read line by line from remote server's files
* Support filename include wildcard(*?)

## How to Build

Ftpbeat uses Glide for dependency management. To install glide see: https://github.com/Masterminds/glide

```shell
$ glide update --no-recursive
$ make 
```

## Configuration

Edit mysqlbeat configuration in ```ftpbeat.yml``` .
You can:
 * Define Username/Password to connect to the FTP server
 * Define Current Directory To Get
 * Define Remote Directory To Read
 * Filenames by using String Array

## How to use
Just run ```ftpbeat -c ftpbeat.yml``` and you are good to go.

## License
GNU General Public License v2
