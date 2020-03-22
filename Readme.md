# GOIP Prometheus exporter

Exporter metrics


|Metric name   |Type   |Meaning   |
|---|---|---|
|goip_up|counter|1 if exporter reached GOIP instance, 0 if not|
|goip_gsm_sim|counter|Sim card 1: detected 0: No Sim|
|goip_gsm_status|counter|GSM Registration 1: registered 0: No|
|goip_status_line|counter|VOIP Registration 1: registered 0: No|
|goip_line_state|counter|GSM Line status: 1:IDLE 0: call in progress|
|goip_acd|counter|Average call duration, seconds |
|goip_asr|counter|Average success ratio, percent|
|goip_callt|counter|Total call duration|
|goip_connected_count|counter|Connected calls count|
|goip_total_count|counter|Total call counts|
|goip_gsm_signal|gauge| Strength of the GSM Signal  |
|goip_nocall_t|counter|Idle time since last call, minutes|
|goip_rct|gauge|CDR Start time|


Example:

```
# HELP goip_acd acd counter
# TYPE goip_acd gauge
goip_acd 0
# HELP goip_asr asr counter
# TYPE goip_asr gauge
goip_asr 0
# HELP goip_callt callt counter
# TYPE goip_callt gauge
goip_callt 0
# HELP goip_gsm_signal Whether gsm_signal is active.
# TYPE goip_gsm_signal gauge
goip_gsm_signal 31
# HELP goip_gsm_sim Whether GSM is active.
# TYPE goip_gsm_sim counter
goip_gsm_sim 1
# HELP goip_gsm_status Whether GSM is active.
# TYPE goip_gsm_status counter
goip_gsm_status 1
# HELP goip_module_status_gsm Whether GSM module is active.
# TYPE goip_module_status_gsm counter
goip_module_status_gsm 1
# HELP goip_nocall_t nocall_t counter
# TYPE goip_nocall_t gauge
goip_nocall_t 6921
# HELP goip_rct rct counter
# TYPE goip_rct gauge
goip_rct 1.584493736e+09
# HELP goip_status_line Whether status_line is active.
# TYPE goip_status_line gauge
goip_status_line 1
# HELP goip_up High level exporter status
# TYPE goip_up gauge
goip_up 1
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
```

### dirty development notes

New dependencies to go.mod are added usually
```
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
```

to get better code completion in jetbrains idea, one of the options would be
pointing GOROOT to your go SDK , saying to one from `gimme`, and issuing `go mod vendor`
to get copy of the dependencies in your vendor directory.

Some good notes on go exporter 

https://www.percona.com/sites/default/files/presentations/Writing%20Prometheus%20exporters.pdf

 
