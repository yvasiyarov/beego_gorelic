package beego_gorelic

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	metrics "github.com/yvasiyarov/go-metrics"
	"github.com/yvasiyarov/gorelic"
	"time"
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

	if beego.AppConfig.String("runmode") == "dev" {
		agent.Verbose = true
	}
	if verbose, err := beego.AppConfig.Bool("NewrelicVerbose"); err == nil {
		agent.Verbose = verbose
	}

	agent.NewrelicName = beego.AppConfig.String("appname")
	agent.Run()

	beego.InsertFilter("*", beego.BeforeRouter, InitNewRelicTimer)
	beego.InsertFilter("*", beego.FinishRouter, ReportMetricsToNewrelic)

	beego.Info("NewRelic agent started")
}
