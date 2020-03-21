package main

import (
    "bufio"
    "fmt"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "regexp"
    "strconv"
    "strings"
    "time"
)

func makeTimestamp() int64 {
        return time.Now().UnixNano() / int64(time.Millisecond)
}

func recordMetrics(goip_address string, user string, password string) {
     var contents = ""
     var nsec = strconv.FormatInt(makeTimestamp(),10)
     var metricsUrl = fmt.Sprintf("http://%s:%s@%s/default/en_US/status.xml?type=&ajaxcachebust=%s", user, password, goip_address, nsec)
     resp, err := http.Get(metricsUrl)
     if err != nil {
          log.Fatal(err)
          fmt.Printf("%s", err)
          os.Exit(1)
        }
     if resp.StatusCode == http.StatusOK {
        bodyBytes, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Fatal(err)
        }
         contents = string(bodyBytes)
     }
    var scanner = bufio.NewScanner(strings.NewReader(contents))
    for scanner.Scan() {
        var currentline = scanner.Text()
        r := regexp.MustCompile(`\<(.+)\>(.*)<\/(.+)\>`)
        matches := r.FindAllStringSubmatch(currentline, -1)
        if matches == nil { continue }
        if len(matches) == 1 && len(matches[0])==4 {
           var metricName=matches[0][1]
           var metricValue=matches[0][2]
           if metricValue == "Y" {
               metricValue = "1"
           }
           if metricValue == "N" {
               metricValue = "0"
           }
           switch metricName {
           case "l1_module_status_gsm":
               statusGsm := prometheus.NewGauge(prometheus.GaugeOpts{
                   Namespace: "goip",
                   Name: "module_status_gsm",
                   Help: "Whether GSM module is active.",
               })
               if metricValue == "Y" {
                   statusGsm.Set(1)
               } else {
                   statusGsm.Set(0)
               }
               prometheus.MustRegister(statusGsm)
               break
           case "l1_gsm_sim":
               gsmSim := prometheus.NewGauge(prometheus.GaugeOpts{
                   Namespace: "goip",
                   Name: "gsm_sim",
                   Help: "Whether GSM is active.",
               })
               if metricValue == "Y" {
                   gsmSim.Set(1)
               } else {
                   gsmSim.Set(0)
               }
               prometheus.MustRegister(gsmSim)
               break
           case "l1_gsm_status":
               gsmStatus := prometheus.NewGauge(prometheus.GaugeOpts{
                   Namespace: "goip",
                   Name: "gsm_status",
                   Help: "Whether GSM is active.",
               })
               if metricValue == "Y" {
                   gsmStatus.Set(1)
               } else {
                   gsmStatus.Set(0)
               }
               prometheus.MustRegister(gsmStatus)
               break
           case "l1_status_line":
               statusLine := prometheus.NewGauge(prometheus.GaugeOpts{
                   Namespace: "goip",
                   Name: "status_line",
                   Help: "Whether GSM is active.",
               })
               if metricValue == "Y" {
                   statusLine.Set(1)
               } else {
                   statusLine.Set(0)
               }
               prometheus.MustRegister(statusLine)
               break
           case "l1_line_state":
               lineState := prometheus.NewGauge(prometheus.GaugeOpts{
                   Namespace: "goip",
                   Name: "line_state",
                   Help: "Whether GSM is active.",
               })
               if metricValue == "Y" {
                   lineState.Set(1)
               } else {
                   lineState.Set(0)
               }
               prometheus.MustRegister(lineState)
               break
           case "l1_gsm_signal":
               gsmSignal := prometheus.NewGauge(prometheus.GaugeOpts{
                   Namespace: "goip",
                   Name: "gsm_signal",
                   Help: "Whether GSM is active.",
               })
               var floatMetric, _ =strconv.ParseFloat(metricValue, 8)
               gsmSignal.Set(floatMetric)

               prometheus.MustRegister(gsmSignal)
               break
           }
        }
    }
}

func main() {
        recordMetrics("192.168.1.189", "admin","admin")
        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":9177", nil)
}
