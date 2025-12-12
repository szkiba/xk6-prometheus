# xk6-prometheus

A [k6](https://go.k6.io/k6) extension that implements a Prometheus HTTP exporter for k6 metrics.

## Features

- **Real-time Metrics Export** - Expose k6 metrics via HTTP endpoint during test execution
- **Long-Running Tests** - Perfect for continuous load testing and monitoring with Prometheus
- **All Metric Types** - Supports all k6 metric types: Counter, Gauge, Rate, and Trend
- **Full Label Support** - Preserves all k6 metric tags as Prometheus labels
- **Configurable** - Customizable port, host, namespace, and subsystem
- **Built-in Metrics** - Automatically exports all [k6 built-in metrics](https://k6.io/docs/using-k6/metrics/#built-in-metrics)
- **Custom Metrics** - Works seamlessly with custom k6 metrics ([Counter](https://k6.io/docs/javascript-api/k6-metrics/counter/), [Gauge](https://k6.io/docs/javascript-api/k6-metrics/gauge/), [Rate](https://k6.io/docs/javascript-api/k6-metrics/rate/), [Trend](https://k6.io/docs/javascript-api/k6-metrics/trend/))

## Table of Contents

- [Installation](#installation)
  - [Pre-built Binaries](#pre-built-binaries)
  - [Build from Source](#build-from-source)
- [Usage](#usage)
  - [Quick Start](#quick-start)
  - [Configuration Parameters](#configuration-parameters)
  - [Examples](#examples)
- [Metric Format](#metric-format)
- [Prometheus Configuration](#prometheus-configuration)
- [Use Cases](#use-cases)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Pre-built Binaries

Download pre-built k6 binaries with xk6-prometheus from the [Releases](https://github.com/szkiba/xk6-prometheus/releases/) page.

### Build from Source

You can build the k6 binary on various platforms, each with its requirements. The following shows how to build k6 binary with this extension on GNU/Linux distributions.

#### Prerequisites

- **Go** - Latest version (matching [k6](https://github.com/grafana/k6#build-from-source) and [xk6](https://github.com/grafana/xk6#requirements) requirements)
- **Git** - For cloning the repository
- **xk6** - Extension builder for k6

#### Install Latest Release

1. Install xk6:

   ```bash
   go install go.k6.io/xk6/cmd/xk6@latest
   ```

2. Build k6 with xk6-prometheus:

   ```bash
   xk6 build --with github.com/szkiba/xk6-prometheus@latest
   ```

#### Development Build

For local development and testing:

```bash
git clone https://github.com/szkiba/xk6-prometheus.git
cd xk6-prometheus
xk6 build --with github.com/szkiba/xk6-prometheus@latest=.
```

This forces xk6 to use your local clone instead of fetching from the repository.

## Usage

### Quick Start

Run k6 with the Prometheus output extension. By default, metrics are exposed on `http://localhost:5656/metrics`.

```bash
./k6 run -o prometheus script.js
```

Access the metrics endpoint:

```bash
curl http://localhost:5656/metrics
```

### Configuration Parameters

Configure the exporter using query string parameters:

```bash
k6 run -o 'prometheus=param1=value1&param2=value2' script.js
```

> [!TIP]
> Use quotes around the `--out` parameter to escape `&` characters from the shell.

| Parameter              | Description                                                                                                                                                                                               | Default               |
|------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------|
| `namespace`            | [Prometheus namespace](https://prometheus.io/docs/practices/naming/) for exported metrics                                                                                                                | `""` (empty)         |
| `subsystem`            | [Prometheus subsystem](https://prometheus.io/docs/practices/naming/) for exported metrics                                                                                                                | `""` (empty)         |
| `host`                 | Hostname or IP address for HTTP endpoint (empty = listen on all interfaces)                                                                                                                              | `""` (all)           |
| `port`                 | TCP port for HTTP endpoint                                                                                                                                                                               | `5656`               |
| `usehistogramfortime`  | If set to `'true'` or `'yes'`, sets the metric type for trends to a [histogram](https://prometheus.io/docs/concepts/metric_types/#histogram) instead of a [summary](https://prometheus.io/docs/concepts/metric_types/#summary) | `"no"` (uses summary) |

> [!TIP]
> It's recommended to use `k6` as either `namespace` or `subsystem` to prefix metrics with `k6_`.

### Examples

#### Basic Usage

Default configuration (port 5656, all interfaces):

```bash
./k6 run -o prometheus script.js
```

#### Custom Port

Run on a specific port:

```bash
./k6 run -o 'prometheus=port=9090' script.js
```

#### With Namespace

Add `k6_` prefix to all metrics:

```bash
./k6 run -o 'prometheus=namespace=k6' script.js
```

#### Custom Host and Port

Listen only on localhost with custom port:

```bash
./k6 run -o 'prometheus=host=127.0.0.1&port=8080' script.js
```

#### Long-Running Test

Run a continuous load test for monitoring:

```bash
./k6 run -o prometheus --duration 24h --vus 10 script.js
```

### Sample Test Script

```javascript
import http from "k6/http";
import { sleep } from "k6";
import { Counter, Trend } from "k6/metrics";

// Custom metrics
let myCounter = new Counter("my_counter");
let myTrend = new Trend("my_trend");

export default function () {
  const response = http.get("https://test.k6.io");
  
  myCounter.add(1);
  myTrend.add(response.timings.duration);
  
  sleep(1);
}
```

## Metric Format

The extension exports k6 metrics in Prometheus text format. Metric types are mapped as follows:

| k6 Metric Type | Prometheus Type         | Description                                     |
|----------------|-------------------------|------------------------------------------------|
| Counter        | Counter                 | Cumulative metric that only increases         |
| Gauge          | Gauge                   | Metric that can go up or down                 |
| Rate           | Histogram               | Ratio of non-zero values (exported as 0 or 1) |
| Trend          | Summary (or Histogram)  | Statistical aggregations with quantiles       |

### Metric Labels

All k6 metric tags are preserved as Prometheus labels:

- `scenario` - Test scenario name
- `group` - Test group name
- `method` - HTTP method (for HTTP metrics)
- `status` - HTTP status code (for HTTP metrics)
- `url` - Request URL (for HTTP metrics)
- `name` - Metric name
- Custom tags from your test script

## Sample HTTP Response

Example metrics output with `namespace=k6`:

```prometheus
# HELP k6_data_received The amount of received data
# TYPE k6_data_received counter
k6_data_received{group="",scenario="default",tls_version=""} 538700
# HELP k6_data_sent The amount of data sent
# TYPE k6_data_sent counter
k6_data_sent{group="",scenario="default",tls_version=""} 9430
# HELP k6_http_req_blocked Time spent blocked  before initiating the request
# TYPE k6_http_req_blocked summary
k6_http_req_blocked{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.5"} 0.003216
k6_http_req_blocked{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.9"} 0.00461
k6_http_req_blocked{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.95"} 0.005075
k6_http_req_blocked{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="1"} 149.563383
k6_http_req_blocked_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 149.7171700000001
k6_http_req_blocked_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 46
k6_http_req_blocked{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.5"} 0.002711
k6_http_req_blocked{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.9"} 0.01625
k6_http_req_blocked{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.95"} 0.02094
k6_http_req_blocked{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="1"} 247.720682
k6_http_req_blocked_sum{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 247.94551300000003
k6_http_req_blocked_count{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 46
# HELP k6_http_req_blocked_current Time spent blocked  before initiating the request (current)
# TYPE k6_http_req_blocked_current gauge
k6_http_req_blocked_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 0.005075
k6_http_req_blocked_current{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 0.00299
# HELP k6_http_req_connecting Time spent establishing TCP connection
# TYPE k6_http_req_connecting summary
k6_http_req_connecting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.5"} 0
k6_http_req_connecting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.9"} 0
k6_http_req_connecting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.95"} 0
k6_http_req_connecting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="1"} 122.939469
k6_http_req_connecting_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 122.939469
k6_http_req_connecting_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 46
k6_http_req_connecting{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.5"} 0
k6_http_req_connecting{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.9"} 0
k6_http_req_connecting{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.95"} 0
k6_http_req_connecting{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="1"} 123.300371
k6_http_req_connecting_sum{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 123.300371
k6_http_req_connecting_count{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 46
# HELP k6_http_req_connecting_current Time spent establishing TCP connection (current)
# TYPE k6_http_req_connecting_current gauge
k6_http_req_connecting_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 0
k6_http_req_connecting_current{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 0
# HELP k6_http_req_duration Total time for the request
# TYPE k6_http_req_duration summary
k6_http_req_duration{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.5"} 122.284411
k6_http_req_duration{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.9"} 122.467859
k6_http_req_duration{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.95"} 122.754156
k6_http_req_duration{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="1"} 125.132449
k6_http_req_duration_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 5624.683091
k6_http_req_duration_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 46
k6_http_req_duration{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.5"} 125.11121
k6_http_req_duration{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.9"} 126.739617
k6_http_req_duration{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.95"} 247.217049
k6_http_req_duration{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="1"} 248.484052
k6_http_req_duration_sum{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 6126.750080999999
k6_http_req_duration_count{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 46
# HELP k6_http_req_duration_current Total time for the request (current)
# TYPE k6_http_req_duration_current gauge
k6_http_req_duration_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 122.316288
k6_http_req_duration_current{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 124.662895
# HELP k6_http_req_failed The rate of failed requests
# TYPE k6_http_req_failed histogram
k6_http_req_failed_bucket{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",le="0"} 46
k6_http_req_failed_bucket{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",le="+Inf"} 46
k6_http_req_failed_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 0
k6_http_req_failed_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 46
k6_http_req_failed_bucket{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",le="0"} 46
k6_http_req_failed_bucket{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",le="+Inf"} 46
k6_http_req_failed_sum{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 0
k6_http_req_failed_count{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 46
# HELP k6_http_req_receiving Time spent receiving response data
# TYPE k6_http_req_receiving summary
k6_http_req_receiving{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.5"} 0.081936
k6_http_req_receiving{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.9"} 0.104153
k6_http_req_receiving{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.95"} 0.112347
k6_http_req_receiving{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="1"} 0.11385
k6_http_req_receiving_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 3.6835980000000004
k6_http_req_receiving_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 46
k6_http_req_receiving{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.5"} 0.083326
k6_http_req_receiving{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.9"} 0.21554
k6_http_req_receiving{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.95"} 122.142645
k6_http_req_receiving{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="1"} 122.322618
k6_http_req_receiving_sum{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 370.93958899999996
k6_http_req_receiving_count{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 46
# HELP k6_http_req_receiving_current Time spent receiving response data (current)
# TYPE k6_http_req_receiving_current gauge
k6_http_req_receiving_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 0.068918
k6_http_req_receiving_current{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 0.086082
# HELP k6_http_req_sending Time spent sending data
# TYPE k6_http_req_sending summary
k6_http_req_sending{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.5"} 0.012458
k6_http_req_sending{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.9"} 0.030467
k6_http_req_sending{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.95"} 0.035746
k6_http_req_sending{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="1"} 0.126034
k6_http_req_sending_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 0.8089879999999999
k6_http_req_sending_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 46
k6_http_req_sending{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.5"} 0.012556
k6_http_req_sending{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.9"} 0.01823
k6_http_req_sending{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.95"} 0.026959
k6_http_req_sending{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="1"} 0.066358
k6_http_req_sending_sum{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 0.689001
k6_http_req_sending_count{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 46
# HELP k6_http_req_sending_current Time spent sending data (current)
# TYPE k6_http_req_sending_current gauge
k6_http_req_sending_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 0.01234
k6_http_req_sending_current{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 0.012584
# HELP k6_http_req_tls_handshaking Time spent handshaking TLS session
# TYPE k6_http_req_tls_handshaking summary
k6_http_req_tls_handshaking{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.5"} 0
k6_http_req_tls_handshaking{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.9"} 0
k6_http_req_tls_handshaking{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.95"} 0
k6_http_req_tls_handshaking{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="1"} 0
k6_http_req_tls_handshaking_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 0
k6_http_req_tls_handshaking_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 46
k6_http_req_tls_handshaking{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.5"} 0
k6_http_req_tls_handshaking{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.9"} 0
k6_http_req_tls_handshaking{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.95"} 0
k6_http_req_tls_handshaking{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="1"} 124.308137
k6_http_req_tls_handshaking_sum{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 124.308137
k6_http_req_tls_handshaking_count{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 46
# HELP k6_http_req_tls_handshaking_current Time spent handshaking TLS session (current)
# TYPE k6_http_req_tls_handshaking_current gauge
k6_http_req_tls_handshaking_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 0
k6_http_req_tls_handshaking_current{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 0
# HELP k6_http_req_waiting Time spent waiting for response
# TYPE k6_http_req_waiting summary
k6_http_req_waiting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.5"} 122.205979
k6_http_req_waiting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.9"} 122.381221
k6_http_req_waiting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="0.95"} 122.639074
k6_http_req_waiting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io",quantile="1"} 124.892565
k6_http_req_waiting_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 5620.190505
k6_http_req_waiting_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 46
k6_http_req_waiting{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.5"} 124.969742
k6_http_req_waiting{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.9"} 126.042663
k6_http_req_waiting{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="0.95"} 126.218584
k6_http_req_waiting{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/",quantile="1"} 126.747568
k6_http_req_waiting_sum{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 5755.121490999999
k6_http_req_waiting_count{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 46
# HELP k6_http_req_waiting_current Time spent waiting for response (current)
# TYPE k6_http_req_waiting_current gauge
k6_http_req_waiting_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 122.23503
k6_http_req_waiting_current{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 124.564229
# HELP k6_http_reqs How many HTTP requests has k6 generated, in total
# TYPE k6_http_reqs counter
k6_http_reqs{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="308",tls_version="",url="http://test.k6.io"} 46
k6_http_reqs{expected_response="true",group="",method="GET",name="https://test.k6.io/",proto="HTTP/1.1",scenario="default",status="200",tls_version="tls1.3",url="https://test.k6.io/"} 46
# HELP k6_iteration_duration The time it took to complete one full iteration
# TYPE k6_iteration_duration summary
k6_iteration_duration{group="",scenario="default",tls_version="",quantile="0.5"} 1248.52603
k6_iteration_duration{group="",scenario="default",tls_version="",quantile="0.9"} 1249.698125
k6_iteration_duration{group="",scenario="default",tls_version="",quantile="0.95"} 1370.836179
k6_iteration_duration{group="",scenario="default",tls_version="",quantile="1"} 1650.467963
k6_iteration_duration_sum{group="",scenario="default",tls_version=""} 56939.48187700001
k6_iteration_duration_count{group="",scenario="default",tls_version=""} 45
# HELP k6_iteration_duration_current The time it took to complete one full iteration (current)
# TYPE k6_iteration_duration_current gauge
k6_iteration_duration_current{group="",scenario="default",tls_version=""} 1247.360889
# HELP k6_iterations The aggregate number of times the VUs in the test have executed
# TYPE k6_iterations counter
k6_iterations{group="",scenario="default",tls_version=""} 45
# HELP k6_vus Current number of active virtual users
# TYPE k6_vus gauge
k6_vus{tls_version=""} 1
# HELP k6_vus_max Max possible number of virtual users
# TYPE k6_vus_max gauge
k6_vus_max{tls_version=""} 1
```

## Prometheus Configuration

Add the k6 endpoint as a scrape target in your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'k6'
    static_configs:
      - targets: ['localhost:5656']
    scrape_interval: 5s  # Adjust based on your needs
```

For Kubernetes deployments, use service discovery:

```yaml
scrape_configs:
  - job_name: 'k6'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app]
        action: keep
        regex: k6
      - source_labels: [__meta_kubernetes_pod_container_port_number]
        action: keep
        regex: "5656"
```

## Use Cases

### Continuous Load Testing

Run k6 continuously and monitor performance over time with Prometheus and Grafana:

```bash
./k6 run -o prometheus --duration 0 --vus 50 script.js
```

### Integration Testing

Monitor application performance during integration tests:

```bash
./k6 run -o prometheus --iterations 1000 --vus 10 integration-test.js
```

### Capacity Planning

Gradually increase load and observe system behavior:

```bash
./k6 run -o prometheus --stages 10m:100,20m:200,10m:300 script.js
```

### Multi-Service Monitoring

Run multiple k6 instances with different ports to test different services:

```bash
# Service A
./k6 run -o 'prometheus=port=5656' service-a-test.js &

# Service B  
./k6 run -o 'prometheus=port=5657' service-b-test.js &
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
