# prometheus-client-go

The Prometheus Go client provides:

1. Built-in Go metrics (memory usage, goroutines, GC, …)
2. The ability to create custom metrics
3. An HTTP handler for the /metrics endpoint

---

## What is a label?
The idea of a metric seems fairly simple. However, there is a problem with such simplicity. On the diagram above, Prometheus monitors several application servers simultaneously. Each of these servers reports mem_used_bytes metric. At any given time, how should Prometheus store multiple samples behind a single metric name?

The first option is aggregation. Prometheus could sum up all the bytes and store the total memory usage of the whole fleet. Or compute an average memory usage and store it. Or min/max memory usage. Or compute and store all of those together. However, there is always a problem with storing only aggregated metrics - we wouldn't be able to pin down a particular server with a bizarre memory usage pattern based on such data.

Luckily, Prometheus uses another approach - it can differentiate samples with the same metric name by labeling them. A label is a certain attribute of a metric. Generally, labels are populated by metric producers (servers in the example above). However, in Prometheus, it's possible to enrich a metric with some static labels based on the producer's identity while recording it on the Prometheus node's side. In the wild, it's common for a Prometheus metric to carry multiple labels.

Typical examples of labels are:

- instance - an instance (a server or cronjob process) of a job being monitored in the <host>:<port> form
- job - a name of a logical group of instances sharing the same purpose
- endpoint - name of an HTTP API endpoint
- method - HTTP method
- status_code - HTTP status code


### Notation
Given a metric name and a set of labels, time series are frequently identified using this notation:
```
<metric name>{<label name>=<label value>, ...}
```
For example, a time series with the metric name api_http_requests_total and the labels method="POST" and handler="/messages" could be written like this:
```
api_http_requests_total{method="POST", handler="/messages"}
```

### Aggregating labels
If a label is not specified, the result of a query will return as many time-series as there are combinations of labels and label values. In order to collapse these combinations (partially or fully), aggregation operators can be used:
```
# Total number of requests regardless of any labels
sum(http_requests_total)

# Average memory usage for each service
avg(memory_used_bytes) by (service)

# Startup time of the oldest instance of each service in each datacenter
min(startup_time_milliseconds) by (service, datacenter)
```

## Caution: 

Remember that every unique combination of key-value label pairs represents a new time series, which can dramatically increase the amount of data stored. Do not use labels to store dimensions with high cardinality (many different label values), such as user IDs, email addresses, or other unbounded sets of values. 

Keeping in mind that each unique combination of labels is considered as a separate time-series, it is usually good practice to keep the number of labels small.

Despite being born in the age of distributed systems, every Prometheus server node is autonomous. I.e., there is no distributed metric storage in the default Prometheus setup, and every node acts as a self-sufficient monitoring server with local metric storage. It simplifies a lot of things, including the following explanation, because we don't need to think of how to merge overlapping series from different Prometheus nodes