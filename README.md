# xk6-prometheus

A k6 extension implements Prometheus HTTP exporter as k6 output extension.

Using **xk6-prometheus** output extension you can collect metrics from long running k6 process with Prometheus. All custom k6 metrics ([Counter](https://k6.io/docs/javascript-api/k6-metrics/counter/),[Gauge](https://k6.io/docs/javascript-api/k6-metrics/gauge/),[Rate](https://k6.io/docs/javascript-api/k6-metrics/rate/),[Trend](https://k6.io/docs/javascript-api/k6-metrics/trend/)) and [build-in metrics](https://k6.io/docs/using-k6/metrics/#built-in-metrics) will be accessible as appropiate Prometheus metrics on a given HTTP port in Prometheus importable text format. 

Built for [k6](https://go.k6.io/k6) using [xk6](https://github.com/grafana/xk6).

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download `xk6`:
  ```bash
  $ go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```bash
  $ xk6 build --with github.com/szkiba/xk6-prometheus@latest
  ```

> You should use at least `v0.31.0` version because xk6-prometheus extension registers itself as output extension. This feature introduced in the `v0.31.0` version of k6.

## Usage

### With defaults

Without parameters the Prometheus HTTP exporter will accessible on port `5656`.

```plain
$ ./k6 run -d 1m --out prometheus script.js

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: script.js
     output: prometheus (:5656)

  scenarios: (100.00%) 1 scenario, 1 max VUs, 1m30s max duration (incl. graceful stop):
           * default: 1 looping VUs for 1m0s (gracefulStop: 30s)


running (1m01.0s), 0/1 VUs, 54 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  1m0s

     data_received..................: 611 kB 10 kB/s
     data_sent......................: 4.1 kB 67 B/s
     http_req_blocked...............: avg=3.37ms   min=2.86µs   med=3.82µs   max=181.96ms p(90)=11.15µs  p(95)=13.52µs 
     http_req_connecting............: avg=2.19ms   min=0s       med=0s       max=118.34ms p(90)=0s       p(95)=0s      
     http_req_duration..............: avg=125.14ms min=118.99ms med=120.68ms max=237.66ms p(90)=121.45ms p(95)=124.07ms
       { expected_response:true }...: avg=125.14ms min=118.99ms med=120.68ms max=237.66ms p(90)=121.45ms p(95)=124.07ms
     http_req_failed................: 0.00%  ✓ 0   ✗ 54 
     http_req_receiving.............: avg=5.1ms    min=85.32µs  med=792.2µs  max=118.29ms p(90)=860.41µs p(95)=903.71µs
     http_req_sending...............: avg=20.68µs  min=12.53µs  med=16.69µs  max=75.97µs  p(90)=29.39µs  p(95)=37.87µs 
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s       p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=120.01ms min=118.17ms med=119.78ms max=127.48ms p(90)=120.6ms  p(95)=120.71ms
     http_reqs......................: 54     0.885451/s
     iteration_duration.............: avg=1.12s    min=1.11s    med=1.12s    max=1.3s     p(90)=1.12s    p(95)=1.16s   
     iterations.....................: 54     0.885451/s
     vus............................: 1      min=1 max=1
     vus_max........................: 1      min=1 max=1
```

### Parameters

The output extension accept parameters in a standard query string format:

```
k6 run --out 'prometheus=param1=value1&param2=value2&param3=value3'
```

> Note apostrophe (`'`) characters around the `--out` parameter! You should use it for escape `&` characters from shell (or use backslash before `&` characters).

The following paremeters are recognized:

parameter | description
----------|------------
namespace | [Prometheus namespace](https://prometheus.io/docs/practices/naming/) for exported metrics (default: "", empty)
subsystem | [Prometheus subsystem](https://prometheus.io/docs/practices/naming/) for exported metrics (default: "", empty)
host      | Hostname or IP address for HTTP endpoint (default: "", empty, listen on all interfaces)
port      | TCP port for HTTP endoint (default: 5656)

*It is recommended to use `k6` as either `namespace` or `subsystem` to prefix exported metrics names with `k6_` string.*

## Sample HTTP response

Here is the relevant part of the metrics HTTP response:

```plain
# HELP k6_data_received The amount of received data
# TYPE k6_data_received counter
k6_data_received{group="",scenario="default"} 430237
# HELP k6_data_sent The amount of data sent
# TYPE k6_data_sent counter
k6_data_sent{group="",scenario="default"} 2888
# HELP k6_http_req_blocked Time spent blocked  before initiating the request
# TYPE k6_http_req_blocked summary
k6_http_req_blocked{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.5"} 0.00377
k6_http_req_blocked{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.9"} 0.008411
k6_http_req_blocked{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.95"} 0.016064
k6_http_req_blocked{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="1"} 872.706547
k6_http_req_blocked_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 872.8805460000001
k6_http_req_blocked_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 38
# HELP k6_http_req_blocked_current Time spent blocked  before initiating the request (current)
# TYPE k6_http_req_blocked_current gauge
k6_http_req_blocked_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 0.003246
# HELP k6_http_req_connecting Time spent establishing TCP connection
# TYPE k6_http_req_connecting summary
k6_http_req_connecting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.5"} 0
k6_http_req_connecting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.9"} 0
k6_http_req_connecting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.95"} 0
k6_http_req_connecting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="1"} 424.67631
k6_http_req_connecting_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 424.67631
k6_http_req_connecting_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 38
# HELP k6_http_req_connecting_current Time spent establishing TCP connection (current)
# TYPE k6_http_req_connecting_current gauge
k6_http_req_connecting_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 0
# HELP k6_http_req_duration Total time for the request
# TYPE k6_http_req_duration summary
k6_http_req_duration{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.5"} 433.957267
k6_http_req_duration{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.9"} 850.686511
k6_http_req_duration{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.95"} 975.504212
k6_http_req_duration{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="1"} 3292.618482
k6_http_req_duration_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 21587.420325
k6_http_req_duration_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 38
# HELP k6_http_req_duration_current Total time for the request (current)
# TYPE k6_http_req_duration_current gauge
k6_http_req_duration_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 418.772552
# HELP k6_http_req_failed The rate of failed requests
# TYPE k6_http_req_failed histogram
k6_http_req_failed_bucket{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",le="0"} 38
k6_http_req_failed_bucket{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",le="+Inf"} 38
k6_http_req_failed_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 0
k6_http_req_failed_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 38
# HELP k6_http_req_receiving Time spent receiving response data
# TYPE k6_http_req_receiving summary
k6_http_req_receiving{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.5"} 1.729595
k6_http_req_receiving{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.9"} 5.581774
k6_http_req_receiving{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.95"} 512.29381
k6_http_req_receiving{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="1"} 1533.230124
k6_http_req_receiving_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 2532.616402999998
k6_http_req_receiving_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 38
# HELP k6_http_req_receiving_current Time spent receiving response data (current)
# TYPE k6_http_req_receiving_current gauge
k6_http_req_receiving_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 2.20869
# HELP k6_http_req_sending Time spent sending data
# TYPE k6_http_req_sending summary
k6_http_req_sending{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.5"} 0.015493
k6_http_req_sending{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.9"} 0.026709
k6_http_req_sending{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.95"} 0.030378
k6_http_req_sending{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="1"} 0.095429
k6_http_req_sending_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 0.7408489999999998
k6_http_req_sending_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 38
# HELP k6_http_req_sending_current Time spent sending data (current)
# TYPE k6_http_req_sending_current gauge
k6_http_req_sending_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 0.013392
# HELP k6_http_req_tls_handshaking Time spent handshaking TLS session
# TYPE k6_http_req_tls_handshaking summary
k6_http_req_tls_handshaking{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.5"} 0
k6_http_req_tls_handshaking{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.9"} 0
k6_http_req_tls_handshaking{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.95"} 0
k6_http_req_tls_handshaking{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="1"} 0
k6_http_req_tls_handshaking_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 0
k6_http_req_tls_handshaking_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 38
# HELP k6_http_req_tls_handshaking_current Time spent handshaking TLS session (current)
# TYPE k6_http_req_tls_handshaking_current gauge
k6_http_req_tls_handshaking_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 0
# HELP k6_http_req_waiting Time spent waiting for response
# TYPE k6_http_req_waiting summary
k6_http_req_waiting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.5"} 431.363566
k6_http_req_waiting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.9"} 535.981583
k6_http_req_waiting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="0.95"} 975.42239
k6_http_req_waiting{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io",quantile="1"} 1759.361649
k6_http_req_waiting_sum{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 19054.063072999994
k6_http_req_waiting_count{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 38
# HELP k6_http_req_waiting_current Time spent waiting for response (current)
# TYPE k6_http_req_waiting_current gauge
k6_http_req_waiting_current{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 416.55047
# HELP k6_http_reqs How many HTTP requests has k6 generated, in total
# TYPE k6_http_reqs counter
k6_http_reqs{expected_response="true",group="",method="GET",name="http://test.k6.io",proto="HTTP/1.1",scenario="default",status="200",url="http://test.k6.io"} 38
# HELP k6_iteration_duration The time it took to complete one full iteration
# TYPE k6_iteration_duration summary
k6_iteration_duration{group="",scenario="default",quantile="0.5"} 1434.73512
k6_iteration_duration{group="",scenario="default",quantile="0.9"} 1945.377781
k6_iteration_duration{group="",scenario="default",quantile="0.95"} 2301.369932
k6_iteration_duration{group="",scenario="default",quantile="1"} 4293.211801
k6_iteration_duration_sum{group="",scenario="default"} 60485.744295
k6_iteration_duration_count{group="",scenario="default"} 38
# HELP k6_iteration_duration_current The time it took to complete one full iteration (current)
# TYPE k6_iteration_duration_current gauge
k6_iteration_duration_current{group="",scenario="default"} 1419.901911
# HELP k6_iterations The aggregate number of times the VUs in the test have executed
# TYPE k6_iterations counter
k6_iterations{group="",scenario="default"} 38
# HELP k6_vus Current number of active virtual users
# TYPE k6_vus gauge
k6_vus 1
# HELP k6_vus_max Max possible number of virtual users
# TYPE k6_vus_max gauge
k6_vus_max 1
```
