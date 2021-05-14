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


running (1m00.8s), 0/1 VUs, 54 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  1m0s

     data_received..................: 611 kB 10 kB/s
     data_sent......................: 4.1 kB 67 B/s
     http_req_blocked...............: avg=2.55ms   min=4.05µs   med=8.92µs   max=137.67ms p(90)=12.7µs   p(95)=13.15µs 
     http_req_connecting............: avg=2.19ms   min=0s       med=0s       max=118.73ms p(90)=0s       p(95)=0s      
     http_req_duration..............: avg=122.13ms min=119.31ms med=121.19ms max=139.75ms p(90)=123.43ms p(95)=126.27ms
       { expected_response:true }...: avg=122.13ms min=119.31ms med=121.19ms max=139.75ms p(90)=123.43ms p(95)=126.27ms
     http_req_failed................: 0.00%  ✓ 0   ✗ 54 
     http_req_receiving.............: avg=839.36µs min=610.96µs med=789.06µs max=3.89ms   p(90)=847.57µs p(95)=881.74µs
     http_req_sending...............: avg=46.16µs  min=14.5µs   med=40.56µs  max=177.74µs p(90)=60.48µs  p(95)=66.41µs 
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s       p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=121.24ms min=118.48ms med=120.31ms max=137.43ms p(90)=122.64ms p(95)=125.43ms
     http_reqs......................: 54     0.888369/s
     iteration_duration.............: avg=1.12s    min=1.11s    med=1.12s    max=1.25s    p(90)=1.12s    p(95)=1.13s   
     iterations.....................: 54     0.888369/s
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
k6_data_received 588452
# HELP k6_data_sent The amount of data sent
# TYPE k6_data_sent counter
k6_data_sent 3952
# HELP k6_http_req_blocked Time spent blocked  before initiating the request
# TYPE k6_http_req_blocked summary
k6_http_req_blocked{quantile="0.5"} 0.008816
k6_http_req_blocked{quantile="0.9"} 0.012707
k6_http_req_blocked{quantile="0.95"} 0.012851
k6_http_req_blocked_sum 138.14527600000005
k6_http_req_blocked_count 52
# HELP k6_http_req_blocked_value Time spent blocked  before initiating the request (value)
# TYPE k6_http_req_blocked_value gauge
k6_http_req_blocked_value 0.009559
# HELP k6_http_req_connecting Time spent establishing TCP connection
# TYPE k6_http_req_connecting summary
k6_http_req_connecting{quantile="0.5"} 0
k6_http_req_connecting{quantile="0.9"} 0
k6_http_req_connecting{quantile="0.95"} 0
k6_http_req_connecting_sum 118.737106
k6_http_req_connecting_count 52
# HELP k6_http_req_connecting_value Time spent establishing TCP connection (value)
# TYPE k6_http_req_connecting_value gauge
k6_http_req_connecting_value 0
# HELP k6_http_req_duration Total time for the request
# TYPE k6_http_req_duration summary
k6_http_req_duration{quantile="0.5"} 121.216566
k6_http_req_duration{quantile="0.9"} 123.535732
k6_http_req_duration{quantile="0.95"} 127.755531
k6_http_req_duration_sum 6354.830187
k6_http_req_duration_count 52
# HELP k6_http_req_duration_value Total time for the request (value)
# TYPE k6_http_req_duration_value gauge
k6_http_req_duration_value 120.766806
# HELP k6_http_req_failed The rate of failed requests
# TYPE k6_http_req_failed histogram
k6_http_req_failed_bucket{le="0"} 52
k6_http_req_failed_bucket{le="+Inf"} 52
k6_http_req_failed_sum 0
k6_http_req_failed_count 52
# HELP k6_http_req_receiving Time spent receiving response data
# TYPE k6_http_req_receiving summary
k6_http_req_receiving{quantile="0.5"} 0.787394
k6_http_req_receiving{quantile="0.9"} 0.847776
k6_http_req_receiving{quantile="0.95"} 0.882193
k6_http_req_receiving_sum 43.79441599999999
k6_http_req_receiving_count 52
# HELP k6_http_req_receiving_value Time spent receiving response data (value)
# TYPE k6_http_req_receiving_value gauge
k6_http_req_receiving_value 0.725954
# HELP k6_http_req_sending Time spent sending data
# TYPE k6_http_req_sending summary
k6_http_req_sending{quantile="0.5"} 0.040285
k6_http_req_sending{quantile="0.9"} 0.060527
k6_http_req_sending{quantile="0.95"} 0.067514
k6_http_req_sending_sum 2.394856
k6_http_req_sending_count 52
# HELP k6_http_req_sending_value Time spent sending data (value)
# TYPE k6_http_req_sending_value gauge
k6_http_req_sending_value 0.040846
# HELP k6_http_req_tls_handshaking Time spent handshaking TLS session
# TYPE k6_http_req_tls_handshaking summary
k6_http_req_tls_handshaking{quantile="0.5"} 0
k6_http_req_tls_handshaking{quantile="0.9"} 0
k6_http_req_tls_handshaking{quantile="0.95"} 0
k6_http_req_tls_handshaking_sum 0
k6_http_req_tls_handshaking_count 52
# HELP k6_http_req_tls_handshaking_value Time spent handshaking TLS session (value)
# TYPE k6_http_req_tls_handshaking_value gauge
k6_http_req_tls_handshaking_value 0
# HELP k6_http_req_waiting Time spent waiting for response
# TYPE k6_http_req_waiting summary
k6_http_req_waiting{quantile="0.5"} 120.329931
k6_http_req_waiting{quantile="0.9"} 122.741878
k6_http_req_waiting{quantile="0.95"} 126.89048
k6_http_req_waiting_sum 6308.640915000001
k6_http_req_waiting_count 52
# HELP k6_http_req_waiting_value Time spent waiting for response (value)
# TYPE k6_http_req_waiting_value gauge
k6_http_req_waiting_value 120.000006
# HELP k6_http_reqs How many HTTP requests has k6 generated, in total
# TYPE k6_http_reqs counter
k6_http_reqs 52
# HELP k6_iteration_duration The time it took to complete one full iteration
# TYPE k6_iteration_duration summary
k6_iteration_duration{quantile="0.5"} 1122.287172
k6_iteration_duration{quantile="0.9"} 1125.321538
k6_iteration_duration{quantile="0.95"} 1139.260684
k6_iteration_duration_sum 58539.15245099999
k6_iteration_duration_count 52
# HELP k6_iteration_duration_value The time it took to complete one full iteration (value)
# TYPE k6_iteration_duration_value gauge
k6_iteration_duration_value 1121.81436
# HELP k6_iterations The aggregate number of times the VUs in the test have executed
# TYPE k6_iterations counter
k6_iterations 52
# HELP k6_vus Current number of active virtual users
# TYPE k6_vus gauge
k6_vus 1
# HELP k6_vus_max Max possible number of virtual users
# TYPE k6_vus_max gauge
k6_vus_max 1
```
