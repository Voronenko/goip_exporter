

### dirty development notes

New dependencies to go.mod are added usually
```
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
```

to get better code completion in jetbrains idea, one of the options would be
pointing GOROOT to your go SDK , saying to one from `gimme`, and issuing `go mod vendor`
to get copy of the dependencies in your vendor directory 
