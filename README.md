# xk6-prometheus

A k6 extension implements Prometheus HTTP exporter as k6 output extension.

Using **xk6-prometheus** output extension you can collect metrics from long running k6 process with Prometheus. All custom k6 metrics ([Counter](https://k6.io/docs/javascript-api/k6-metrics/counter/),[Gauge](https://k6.io/docs/javascript-api/k6-metrics/gauge/),[Rate](https://k6.io/docs/javascript-api/k6-metrics/rate/),[Trend](https://k6.io/docs/javascript-api/k6-metrics/trend/)) and [build-in metrics](https://k6.io/docs/using-k6/metrics/#built-in-metrics) will be accessible as appropiate Prometheus metrics on a given HTTP port in Prometheus importable text format. 

Built for [k6](https://go.k6.io/k6) using [xk6](https://github.com/grafana/xk6).

## Download

You can download pre-built k6 binaries from [Releases](https://github.com/szkiba/xk6-prometheus/releases/) page. Check [Packages](https://github.com/szkiba/xk6-prometheus/pkgs/container/xk6-prometheus) page for pre-built k6 Docker images.

## Build

You can build the k6 binary on various platforms, each with its requirements. The following shows how to build k6 binary with this extension on GNU/Linux distributions.

### Prerequisites

You must have the latest Go version installed to build the k6 binary. The latest version should match [k6](https://github.com/grafana/k6#build-from-source) and [xk6](https://github.com/grafana/xk6#requirements).

- [Git](https://git-scm.com/) for cloning the project
- [xk6](https://github.com/grafana/xk6) for building k6 binary with extensions

### Install and build the latest tagged version

1. Install `xk6`:

   ```shell
   go install go.k6.io/xk6/cmd/xk6@latest
   ```

2. Build the binary:

   ```shell
   xk6 build --with github.com/szkiba/xk6-prometheus@latest
   ```

> **Note**
> You can always use the latest version of k6 to build the extension, but the earliest version of k6 that supports extensions via xk6 is v0.43.0. The xk6 is constantly evolving, so some APIs may not be backward compatible.

### Build for development

If you want to add a feature or make a fix, clone the project and build it using the following commands. The xk6 will force the build to use the local clone instead of fetching the latest version from the repository. This process enables you to update the code and test it locally.

```bash
git clone git@github.com:szkiba/xk6-prometheus.git && cd xk6-prometheus
xk6 build --with github.com/szkiba/xk6-prometheus@latest=.
```

## Docker

You can also use pre-built k6 image within a Docker container. In order to do that, you will need to execute something like the following:

**Linux**

```plain
docker run -v $(pwd):/scripts -it --rm ghcr.io/szkiba/xk6-prometheus:latest run -d 1m --out=dashboard /scripts/script.js
```

**Windows**

```plain
docker run -v %cd%:/scripts -it --rm ghcr.io/szkiba/xk6-prometheus:latest run -d 1m --out=prometheus /scripts/script.js
```

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
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 2.4485e-05
go_gc_duration_seconds{quantile="0.25"} 4.5041e-05
go_gc_duration_seconds{quantile="0.5"} 5.3342e-05
go_gc_duration_seconds{quantile="0.75"} 5.9793e-05
go_gc_duration_seconds{quantile="1"} 0.000151794
go_gc_duration_seconds_sum 0.000642498
go_gc_duration_seconds_count 11
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 24
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.18.1"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 4.038792e+07
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 1.40677416e+08
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 5496
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 2.013702e+06
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 6.310184e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 4.038792e+07
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 4.907008e+06
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 5.2830208e+07
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 172212
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 2.473984e+06
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 5.7737216e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.6847687042787058e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 2.185914e+06
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 9600
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 657152
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 685440
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 5.153816e+07
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.391104e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 983040
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 983040
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 6.712808e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 13
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
