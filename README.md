# Ftpbeat
Fully customizable Beat for FTPL Server - this beat can ship the results of any query defined on the config file to Elasticsearch.


## Current status
Ftpbeat still on beta.

#### To Do:
* Add support for sFtp


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
 * Define Remote Directory
 * Filenames by using string array

## How to use
Just run ```ftpbeat -c ftpbeat.yml``` and you are good to go.

## License
GNU General Public License v2
