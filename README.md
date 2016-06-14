ssl_certificate_exporter
=============

Exporter for SSL Certificate metrics https://prometheus.io/ 

## Configuration

1. Write the domain you want to check the expiration date of the ssl to JSON File.

```json
{"domains":["example.com","example.net","example.org"]}
```

2. Hosting the JSON File to the web, such as Gist and S3.

3. Specify it in the environment variable CONFIG_URL.

ex. https://gist.githubusercontent.com/s-aska/03c41cf0d3f8b369cf0ae80d02a26c02/raw/3c742b80c4c1c7e79fb6705cda19808efb8048eb/config.json

## Run

### Local machine

```sh
git clone git@github.com:s-aska/ssl_certificate_exporter.git
cd s-aska/ssl_certificate_exporter
go build
env \
    PORT=9100 \
    CONFIG_URL=https://gist.githubusercontent.com/s-aska/03c41cf0d3f8b369cf0ae80d02a26c02/raw/3c742b80c4c1c7e79fb6705cda19808efb8048eb/config.json \
    ./ssl_certificate_exporter
```

### Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

### Docker

```
docker pull aska/ssl_certificate_exporter

docker run -e PORT=9100 \
    -e CONFIG_URL="https://.../config.json" \
    -p 9100:9100 \
    --name ssl_certificate_exporter \
    --rm \
    aska/ssl_certificate_exporter
```

## Reloading configuration

```sh
curl -X POST http://XXX.XXX.XXX.XXX:9100/-/reload
```

## Grafana (example)

```
Query: ssl_certificate_expires
Legend format: {{domain}}
Axes Left Y Unit: seconds(s)
```

![Grafana](https://github.com/s-aska/ssl_certificate_exporter/wiki/grafana-1.png "Grafana")
