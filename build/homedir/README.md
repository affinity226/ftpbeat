# Welcome to ftpbeat 1.2.0

Sends events to Elasticsearch or Logstash

## Getting Started

To get started with ftpbeat, you need to set up Elasticsearch on your localhost first. After that, start ftpbeat with:

     ./ftpbeat  -c ftpbeat.yml -e

This will start the beat and send the data to your Elasticsearch instance. To load the dashboards for ftpbeat into Kibana, run:

    ./scripts/import_dashboards

For further steps visit the [Getting started](https://www.elastic.co/guide/en/beats/ftpbeat/5.2/ftpbeat-getting-started.html) guide.

## Documentation

Visit [Elastic.co Docs](https://www.elastic.co/guide/en/beats/ftpbeat/5.2/index.html) for the full ftpbeat documentation.

## Release notes

https://www.elastic.co/guide/en/beats/libbeat/5.2/release-notes-1.2.0.html
