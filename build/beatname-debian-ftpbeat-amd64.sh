#!/bin/bash

# Script to run ftpbeat in foreground with the same path settings that
# the init script / systemd unit file would do.

/usr/share/ftpbeat/bin/ftpbeat \
  -path.home /usr/share/ftpbeat \
  -path.config /etc/ftpbeat \
  -path.data /var/lib/ftpbeat \
  -path.logs /var/log/ftpbeat \
  $@
