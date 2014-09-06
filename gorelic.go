package beego_gorelic

import (
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	metrics "github.com/yvasiyarov/go-metrics"
	"github.com/yvasiyarov/gorelic"
)

const (
	SEPARATOR = "-"
)

var agent *gorelic.Agent

func InitNewRelicTimer(ctx *context.Context) {
	startTime := time.Now()
	ctx.Input.SetData("newrelic_timer", startTime)
}
func ReportMetricsToNewrelic(ctx *context.Context) {
	startTimeInterface := ctx.Input.GetData("newrelic_timer")
	if startTime, ok := startTimeInterface.(time.Time); ok {
		agent.HTTPTimer.UpdateSince(startTime)
	}
}

func InitNewrelicAgent() {

	license := beego.AppConfig.String("NewrelicLicense")
	if license == "" {
		beego.Warn("Please specify NewRelic license in the application config: NewrelicLicense=7bceac019c7dcafae1ef95be3e3a3ff8866de245")
		return
	}

	agent = gorelic.NewAgent()
	agent.NewrelicLicense = license

	agent.HTTPTimer = metrics.NewTimer()
	agent.CollectHTTPStat = true

	if beego.RunMode == "dev" {
		agent.Verbose = true
	}
	if verbose, err := beego.AppConfig.Bool("NewrelicVerbose"); err == nil {
		agent.Verbose = verbose
	}

	nameParts := []string{beego.AppConfig.String("appname")}
	switch strings.ToLower(beego.AppConfig.String("NewrelicAppnameRunmode")) {
	case "append":
		nameParts = append(nameParts, beego.RunMode)

	case "prepend":
		nameParts = append([]string{beego.RunMode}, nameParts...)
	}
	agent.NewrelicName = strings.Join(nameParts, SEPARATOR)
	agent.Run()

	beego.InsertFilter("*", beego.BeforeRouter, InitNewRelicTimer)
	beego.InsertFilter("*", beego.FinishRouter, ReportMetricsToNewrelic)

	beego.Info("NewRelic agent started")
}
