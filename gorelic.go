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
	var appname string
	license := beego.AppConfig.String("newrelicLicense")
	if license == "" {
		beego.Warn("Please specify NewRelic license in the application config: NewrelicLicense=7bceac019c7dcafae1ef95be3e3a3ff8866de245")
		return
	}

	agent = gorelic.NewAgent()
	agent.NewrelicLicense = license

	agent.HTTPTimer = metrics.NewTimer()
	agent.CollectHTTPStat = true

	if beego.BConfig.RunMode == "dev" {
		agent.Verbose = true
	}
	if verbose, err := beego.AppConfig.Bool("newrelicVerbose"); err == nil {
		agent.Verbose = verbose
	}
	// Checking if New Relic appname overrides the default appname
	appname = beego.AppConfig.String("newrelicAppname")
	if appname == "" {
		// If not set revert to using beego appname as default
		appname = beego.AppConfig.String("appname")
	}
	nameParts := []string{appname}

	switch strings.ToLower(beego.AppConfig.String("newrelicAppnameRunmode")) {
	case "append":
		nameParts = append(nameParts, beego.BConfig.RunMode)

	case "prepend":
		nameParts = append([]string{beego.BConfig.RunMode}, nameParts...)
	}
	agent.NewrelicName = strings.Join(nameParts, SEPARATOR)
	agent.Run()

	beego.InsertFilter("*", beego.BeforeRouter, InitNewRelicTimer, false)
	beego.InsertFilter("*", beego.FinishRouter, ReportMetricsToNewrelic, false)

	beego.Info("NewRelic agent started")
}
