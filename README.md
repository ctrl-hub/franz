# Franz

Franz is a prometheus metrics exporter for Confluent Cloud Kafka instances. It creates a ClusterAdmin connection to the cluster, queries it and exposes the metrics back in prometheus format.

## Why?

Our need to export this data led us down the path of kafka-exporter and it's fork. Neither have good support for env variables, which are a basic part of a 12FA and they do not appear to be maintained, or PRs accepted.

## Metrics

```bash
# HELP kafka_broker_info Broker information - the value is arbitrary (labels hold info).
# TYPE kafka_broker_info gauge
kafka_broker_info{addr="b0-pkc-xxxxxx.europe-west1.gcp.confluent.cloud:9092",cluster="my-cluster",id="0",rack="0"} 1
...
# HELP kafka_consumer_group_members The number of members in a consumer group.
# TYPE kafka_consumer_group_members gauge
kafka_consumer_group_members{cluster="my-cluster",name="svc.email",protocol="",protocol_type="consumer",state="Empty"} 0
...
# HELP kafka_consumer_group_offsets The offset of a consumer group in a topic.
# TYPE kafka_consumer_group_offsets gauge
kafka_consumer_group_offsets{cluster="my-cluster",name="svc.email",partition="0",topic="my-super-topic"} 15648
...
# HELP kafka_topic_info Topic information - the value is arbitrary (labels hold info).
# TYPE kafka_topic_info gauge
kafka_topic_info{cluster="my-cluster",internal="false",name="my-super-topic",partitions="1"} 1
...
# HELP kafka_topic_parititions Number of partitions in a topic
# TYPE kafka_topic_parititions gauge
kafka_topic_parititions{cluster="my-cluster",topic="my-super-topic"} 1
...
# HELP kafka_topic_partition_info Topic partition information - the value is arbitrary (labels hold info).
# TYPE kafka_topic_partition_info gauge
kafka_topic_partition_info{cluster="my-cluster",leader="1",offline_replicas="0",partition="0",replicas="3",topic="my-super-topic"} 1
...
```

## Running Franz

There is an `.env.example` which you should copy to `.env` (`cp .env.example .env`). You should then add your cluster details and credentials to the env vars.

### Configuration Options

The run time configurarion can be changed using either environment variables or by passing arguments to the process.

Env | Flag | Description | Default | Options
--- | ---- | ------- | ----------- | ------- 
`CONFLUENT_CLUSTER_LABEL`   | `--confluent_cluster_label`   | The label value to add for the cluster label in metric
`CONFLUENT_ENDPOINT`        | `--confluent_endpoint`        | The confluent endpoint for the cluster (e.g. `pkc-xxxx.region.provider.onfluent.cloud:9092`)
`CONFLUENT_API_KEY`         | `--confluent_api_key`         | The confluent API Key for SASL authentication (username)
`CONFLUENT_API_SECRET`      | `--confluent_api_secret`      | The confluent API Secret for SASL authentication (password)
`LOG_LEVEL`                 | `--log_level`                 | The log level to use                                                                          | `info`        | `fatal`, `panic`, `warn`, `info`, `debug`
`LOG_FORMAT`                | `--log_format`                | How you want your logs formatted                                                              | `logfmt`      | `logfmt`, `json`
`METRICS_PATH`              | `--metrics_path`              | The path to serve metrics on                                                                  | `/metrics`    |
`METRICS_PORT`              | `--metrics_port`              | The port to serve metrics on                                                                  | `3100`
`POLLING_INTERVAL_SECONDS`  | `--polling_interval_seconds`  | How often to poll the kafka cluster (in seconds)                                              | `10`

### Locally

With Go installed, you can run `go get .` to grab the package deps and then to run from source:

```bash
export $(cat .env | xargs) && go run main.go
```

You can pass in arguments on the command line, for example:

```bash
export $(cat .env | xargs) && go run main.go --polling_interval_seconds 30
```

To compile the app, you can run:

```bash
go build -o bin/franz
```

You can then run `export $(cat .env | xargs) && ./bin/franz --polling_interval_seconds 30` etc

To install the app to your local `$GOPATH` you can run:

```bash
go install
```

You can then run:

```bash
export $(cat .env | xargs) && franz --polling_interval_seconds 30
```

### Docker Compose

```bash
docker-compose up -d --build
```

### Kubernetes

See [k8s/install.yaml](./k8s/install.yaml) for an example of deploying to kubernetes. It is not intended to be rpoduction ready, just to give you an example.

You will need to tailor this to your environment using a ServiceMonitor or manual scrape job in your prometheus config.

You will also want to create your own secret for the sensitive env variables.

### Building your own docker image

```bash
docker build -t franz
```
