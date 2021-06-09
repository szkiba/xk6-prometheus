# xk6-prometheus

A k6 extension implements Prometheus HTTP exporter as k6 output extension.

Using **xk6-prometheus** output extension you can collect metrics from long running k6 process with Prometheus. All custom k6 metrics ([Counter](https://k6.io/docs/javascript-api/k6-metrics/counter/),[Gauge](https://k6.io/docs/javascript-api/k6-metrics/gauge/),[Rate](https://k6.io/docs/javascript-api/k6-metrics/rate/),[Trend](https://k6.io/docs/javascript-api/k6-metrics/trend/)) and [build-in metrics](https://k6.io/docs/using-k6/metrics/#built-in-metrics) will be accessible as appropiate Prometheus metrics on a given HTTP port in Prometheus importable text format. 

Built for [k6](https://go.k6.io/k6) using [xk6](https://github.com/k6io/xk6).

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download `xk6`:
  ```bash
  $ go install github.com/k6io/xk6/cmd/xk6@latest
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
k6_data_received 600077
# HELP k6_data_sent The amount of data sent
# TYPE k6_data_sent counter
k6_data_sent 4028
# HELP k6_http_req_blocked Time spent blocked  before initiating the request
# TYPE k6_http_req_blocked summary
k6_http_req_blocked{quantile="0.5"} 0.003785
k6_http_req_blocked{quantile="0.9"} 0.011186
k6_http_req_blocked{quantile="0.95"} 0.014801
k6_http_req_blocked{quantile="1"} 181.961168
k6_http_req_blocked_sum 182.2282190000001
k6_http_req_blocked_count 54
# HELP k6_http_req_blocked_current Time spent blocked  before initiating the request (current)
# TYPE k6_http_req_blocked_current gauge
k6_http_req_blocked_current 0.004652
# HELP k6_http_req_connecting Time spent establishing TCP connection
# TYPE k6_http_req_connecting summary
k6_http_req_connecting{quantile="0.5"} 0
k6_http_req_connecting{quantile="0.9"} 0
k6_http_req_connecting{quantile="0.95"} 0
k6_http_req_connecting{quantile="1"} 118.345391
k6_http_req_connecting_sum 118.345391
k6_http_req_connecting_count 54
# HELP k6_http_req_connecting_current Time spent establishing TCP connection (current)
# TYPE k6_http_req_connecting_current gauge
k6_http_req_connecting_current 0
# HELP k6_http_req_duration Total time for the request
# TYPE k6_http_req_duration summary
k6_http_req_duration{quantile="0.5"} 120.68482
k6_http_req_duration{quantile="0.9"} 121.461453
k6_http_req_duration{quantile="0.95"} 128.291489
k6_http_req_duration{quantile="1"} 237.666802
k6_http_req_duration_sum 6757.691293999998
k6_http_req_duration_count 54
# HELP k6_http_req_duration_current Total time for the request (current)
# TYPE k6_http_req_duration_current gauge
k6_http_req_duration_current 120.751952
# HELP k6_http_req_failed The rate of failed requests
# TYPE k6_http_req_failed histogram
k6_http_req_failed_bucket{le="0"} 54
k6_http_req_failed_bucket{le="+Inf"} 54
k6_http_req_failed_sum 0
k6_http_req_failed_count 54
# HELP k6_http_req_receiving Time spent receiving response data
# TYPE k6_http_req_receiving summary
k6_http_req_receiving{quantile="0.5"} 0.79092
k6_http_req_receiving{quantile="0.9"} 0.862976
k6_http_req_receiving{quantile="0.95"} 0.973895
k6_http_req_receiving{quantile="1"} 118.292546
k6_http_req_receiving_sum 275.923961
k6_http_req_receiving_count 54
# HELP k6_http_req_receiving_current Time spent receiving response data (current)
# TYPE k6_http_req_receiving_current gauge
k6_http_req_receiving_current 0.706994
# HELP k6_http_req_sending Time spent sending data
# TYPE k6_http_req_sending summary
k6_http_req_sending{quantile="0.5"} 0.016638
k6_http_req_sending{quantile="0.9"} 0.030783
k6_http_req_sending{quantile="0.95"} 0.041126
k6_http_req_sending{quantile="1"} 0.07597
k6_http_req_sending_sum 1.11689
k6_http_req_sending_count 54
# HELP k6_http_req_sending_current Time spent sending data (current)
# TYPE k6_http_req_sending_current gauge
k6_http_req_sending_current 0.018325
# HELP k6_http_req_tls_handshaking Time spent handshaking TLS session
# TYPE k6_http_req_tls_handshaking summary
k6_http_req_tls_handshaking{quantile="0.5"} 0
k6_http_req_tls_handshaking{quantile="0.9"} 0
k6_http_req_tls_handshaking{quantile="0.95"} 0
k6_http_req_tls_handshaking{quantile="1"} 0
k6_http_req_tls_handshaking_sum 0
k6_http_req_tls_handshaking_count 54
# HELP k6_http_req_tls_handshaking_current Time spent handshaking TLS session (current)
# TYPE k6_http_req_tls_handshaking_current gauge
k6_http_req_tls_handshaking_current 0
# HELP k6_http_req_waiting Time spent waiting for response
# TYPE k6_http_req_waiting summary
k6_http_req_waiting{quantile="0.5"} 119.771619
k6_http_req_waiting{quantile="0.9"} 120.604942
k6_http_req_waiting{quantile="0.95"} 120.858167
k6_http_req_waiting{quantile="1"} 127.48803
k6_http_req_waiting_sum 6480.650443
k6_http_req_waiting_count 54
# HELP k6_http_req_waiting_current Time spent waiting for response (current)
# TYPE k6_http_req_waiting_current gauge
k6_http_req_waiting_current 120.026633
# HELP k6_http_reqs How many HTTP requests has k6 generated, in total
# TYPE k6_http_reqs counter
k6_http_reqs 54
# HELP k6_iteration_duration The time it took to complete one full iteration
# TYPE k6_iteration_duration summary
k6_iteration_duration{quantile="0.5"} 1121.458258
k6_iteration_duration{quantile="0.9"} 1122.620529
k6_iteration_duration{quantile="0.95"} 1237.23951
k6_iteration_duration{quantile="1"} 1304.729788
k6_iteration_duration_sum 59861.31288699999
k6_iteration_duration_count 53
# HELP k6_iteration_duration_current The time it took to complete one full iteration (current)
# TYPE k6_iteration_duration_current gauge
k6_iteration_duration_current 1121.360347
# HELP k6_iterations The aggregate number of times the VUs in the test have executed
# TYPE k6_iterations counter
k6_iterations 53
# HELP k6_vus Current number of active virtual users
# TYPE k6_vus gauge
k6_vus 1
# HELP k6_vus_max Max possible number of virtual users
# TYPE k6_vus_max gauge
k6_vus_max 1
```
