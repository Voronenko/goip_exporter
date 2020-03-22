package main

import (
    "bufio"
    "errors"
    "fmt"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "gopkg.in/alecthomas/kingpin.v2"
    "io/ioutil"
    "log"
    "net/http"
    "regexp"
    "strconv"
    "strings"
    "time"
)

var (
    errKeyNotFound = errors.New("key not found")
)

// Timestamp to beat the cache on GOIP device
func makeTimestamp() int64 {
        return time.Now().UnixNano() / int64(time.Millisecond)
}

//parses metric
func (e *Exporter) parseAndNewMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, stats map[string]string, key string, labelValues ...string) error {
    return e.extractValueAndNewMetric(ch, desc, valueType, parse, stats, key, labelValues...)
}

func (e *Exporter) parseBoolAndNewMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, stats map[string]string, key string, labelValues ...string) error {
    return e.extractValueAndNewMetric(ch, desc, valueType, parseBool, stats, key, labelValues...)
}

func (e *Exporter) parseDateTimeAndNewMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, stats map[string]string, key string, labelValues ...string) error {
    return e.extractValueAndNewMetric(ch, desc, valueType, parseDateTime, stats, key, labelValues...)
}

//extracts metric
func (e *Exporter) extractValueAndNewMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, f func(map[string]string, string) (float64, error), stats map[string]string, key string, labelValues ...string) error {
    v, err := f(stats, key)
    if err == errKeyNotFound {
        return nil
    }
    if err != nil {
        return err
    }

    ch <- prometheus.MustNewConstMetric(desc, valueType, v, labelValues...)
    return nil
}

func parse(stats map[string]string, key string) (float64, error) {
    value, ok := stats[key]
    if !ok {
        log.Printf("Key not found: %s", key)
        return 0, errKeyNotFound
    }

    v, err := strconv.ParseFloat(value, 64)
    if err != nil {
        log.Printf("Failed to parse key %s value %s error %s", key, value, err)
        return 0, err
    }
    return v, nil
}

func parseBool(stats map[string]string, key string) (float64, error) {
    value, ok := stats[key]
    if !ok {
        log.Printf("Key not found %s", key)
        return 0, errKeyNotFound
    }

    switch value {
    case "yes":
        return 1, nil
    case "Y":
        return 1, nil
    case "no":
        return 0, nil
    case "N":
        return 0, nil
    default:
        log.Printf("Failed to parse key %s value %s", key, value)
        return 0, errors.New("failed parse a bool value")
    }
}

func parseDateTime(stats map[string]string, key string) (float64, error) {
    value, ok := stats[key]
    if !ok {
        log.Printf("Key not found: %s", key)
        return 0, errKeyNotFound
    }

    layout := "2006-01-02 15:04:05"
    t, err := time.Parse(layout, value)

    if err != nil {
        log.Printf("Datetime failed to parse for key: %s from %s", key, value)
        return 0, err
    }
    return float64(t.Unix()), nil
}

// Exporter collects metrics from a GOIP voip gateway.
type Exporter struct {
    address string
    user string
    pass string

    up                       *prometheus.Desc
    module_status_gsm        *prometheus.Desc
    gsm_sim                  *prometheus.Desc
    gsm_status               *prometheus.Desc
    status_line              *prometheus.Desc
    line_state               *prometheus.Desc
    gsm_signal               *prometheus.Desc
    nocall_t                 *prometheus.Desc
    acd                      *prometheus.Desc
    asr                      *prometheus.Desc
    callt                    *prometheus.Desc
    rct                      *prometheus.Desc
}

// NewExporter returns an initialized exporter.
func NewExporter(goipAddress string, goipUser string, goipPass string) *Exporter {
    var namespace = "goip"
    return &Exporter{
        address: goipAddress,
        user: goipUser,
        pass:  goipPass,
        up: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "up"),
            "High level exporter status",
            nil,
            nil,
        ),
        module_status_gsm: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "module_status_gsm"),
            "Whether GSM module is active.",
            nil,
            nil,
        ),
        gsm_sim: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "gsm_sim"),
            "Whether GSM is active.",
            nil,
            nil,
        ),
        gsm_status: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "gsm_status"),
            "Whether GSM is active.",
            nil,
            nil,
        ),
        status_line: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "status_line"),
            "Whether status_line is active.",
            nil,
            nil,
        ),
        line_state: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "line_state"),
            "Whether line_state is active.",
            nil,
            nil,
        ),
        gsm_signal: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "gsm_signal"),
            "Whether gsm_signal is active.",
            nil,
            nil,
        ),
        nocall_t: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "nocall_t"),
            "nocall_t counter",
            nil,
            nil,
        ),
        acd: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "acd"),
            "acd counter",
            nil,
            nil,
        ),
        asr: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "asr"),
            "asr counter",
            nil,
            nil,
        ),
        callt: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "callt"),
            "callt counter",
            nil,
            nil,
        ),
        rct: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, "", "rct"),
            "rct counter",
            nil,
            nil,
        ),

    }
}

// Describe describes all the metrics exported by the goip exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
    ch <- e.up
    ch <- e.module_status_gsm
    ch <- e.gsm_sim
    ch <- e.gsm_status
    ch <- e.status_line
    ch <- e.line_state
    ch <- e.gsm_signal
    ch <- e.nocall_t
    ch <- e.acd
    ch <- e.asr
    ch <- e.callt
    ch <- e.rct
}

// Collect fetches the statistics from the configured goip server, and
// delivers them as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
    var stats string
    var nsec = strconv.FormatInt(makeTimestamp(), 10)
    var metricsUrl = fmt.Sprintf("http://%s:%s@%s/default/en_US/status.xml?type=&ajaxcachebust=%s", e.user, e.pass, e.address, nsec)
    resp, err := http.Get(metricsUrl)
    if err != nil {
        ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
        log.Printf("Failed to connect to goip server %s", err)
        return
    }
    if resp.StatusCode == http.StatusOK {
        bodyBytes, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
            log.Printf("Failed to connect to goip server %s", err)
            return
        }
        stats = string(bodyBytes)
    }
    up := float64(1)
    var rawMetrics map[string]string = parseMetricsRaw(stats)
    if err := e.parseStats(ch, rawMetrics); err != nil {
        up = 0
    }
    ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, up)
}

func (e *Exporter) parseStats(ch chan<- prometheus.Metric, statsraw map[string]string) error {
    var parseError error
    var err = firstError(
        e.parseAndNewMetric(ch, e.module_status_gsm, prometheus.CounterValue, statsraw, "l1_module_status_gsm"),
        e.parseBoolAndNewMetric(ch, e.gsm_sim, prometheus.CounterValue, statsraw, "l1_gsm_sim"),
        e.parseBoolAndNewMetric(ch, e.gsm_status, prometheus.CounterValue, statsraw, "l1_gsm_status"),
        e.parseBoolAndNewMetric(ch, e.status_line, prometheus.GaugeValue, statsraw, "l1_status_line"),
        // TODO: <l1_line_state>IDLE</l1_line_state>
        e.parseAndNewMetric(ch, e.gsm_signal, prometheus.GaugeValue, statsraw, "l1_gsm_signal"),
        // TODO: l1_gsm_cur_oper
        e.parseAndNewMetric(ch, e.nocall_t, prometheus.GaugeValue, statsraw, "l1_nocall_t"),
        e.parseAndNewMetric(ch, e.acd, prometheus.GaugeValue, statsraw, "l1_acd"),
        e.parseAndNewMetric(ch, e.asr, prometheus.GaugeValue, statsraw, "l1_asr"),
        e.parseAndNewMetric(ch, e.callt, prometheus.GaugeValue, statsraw, "l1_callt"),
        // TODO: <l1_callc>0/7</l1_callc>
        e.parseDateTimeAndNewMetric(ch, e.rct, prometheus.GaugeValue, statsraw, "l1_rct"),
        //TODO: <l1_sms_count>NO LIMIT</l1_sms_count>
    )
    if err != nil {
        parseError = err
    }
    return parseError
}


func parseMetricsRaw(contents string) map[string]string {
    var scanner = bufio.NewScanner(strings.NewReader(contents))
    var result  =  map[string]string{}

    for scanner.Scan() {
        var currentLine = scanner.Text()
        r := regexp.MustCompile(`\<(.+)\>(.*)<\/(.+)\>`)
        matches := r.FindAllStringSubmatch(currentLine, -1)
        if matches == nil {
            continue
        }
        if len(matches) == 1 && len(matches[0]) == 4 {
            var metricName = matches[0][1]
            var metricValue = matches[0][2]
            result[metricName] = metricValue
        }
   }
   return result
}

func firstError(errors ...error) error {
    for _, v := range errors {
        if v != nil {
            return v
        }
    }
    return nil
}

func main() {
    var (
        goipAddress   = kingpin.Flag("goip.address", "GoIP device address").Default("192.168.1.189").String()
        goipUser       = kingpin.Flag("goip.user", "Username to connect to GOIP device.").Default("admin").String()
        goipPassword       = kingpin.Flag("goip.password", "User's password used to connect to GOIP device.").Default("admin").String()
        goipExporterAddr = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9177").String()
        metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
    )

    kingpin.HelpFlag.Short('h')
    kingpin.Parse()
    log.Print("Starting goip_exporter")
    var exporterInstance = NewExporter(*goipAddress, *goipUser, *goipPassword)
    prometheus.MustRegister(exporterInstance)
    prometheus.Unregister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
    prometheus.Unregister(prometheus.NewGoCollector())
    http.Handle(*metricsPath, promhttp.Handler())
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`<html>
             <head><title>GOIP Exporter</title></head>
             <body>
             <h1>GOIP Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
    })

    log.Printf("Listening on address %s", *goipExporterAddr)
    if err := http.ListenAndServe(*goipExporterAddr, nil); err != nil {
        log.Fatal("Error running HTTP server %s", err )
    }
}