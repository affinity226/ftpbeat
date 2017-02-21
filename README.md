# Ftpbeat
Fully customizable Beat for FTPL Server - this beat can ship the results of any query defined on the config file to Elasticsearch.


## Current status
Ftpbeat still on beta.

#### To Do:
* Add support for sFtp


## Features

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

## How to use
Just run ```ftpbeat -c ftpbeat.yml``` and you are good to go.

## License
GNU General Public License v2
