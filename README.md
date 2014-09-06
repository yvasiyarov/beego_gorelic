beego_gorelic
=============

beego_gorelic is NewRelic middleware for beego framework.

# Installation
 - Run "go get github.com/yvasiyarov/beego_gorelic"
 - Open routers/router.go file and:
     - Add  "github.com/yvasiyarov/beego_gorelic" to import statement
     - Add beego_gorelic.InitNewrelicAgent() call to your init() function
 - Add NewrelicLicense key to conf/app.conf 
 - Optionally add NewrelicVerbose=true if you wanna to see metrics, reported by NewRelic Agent

If your application use runmode=dev, then NewrelicVerbose will be set to true by default

## Optional Configuration
The following key can be added to the conf/app.conf to prepend or append the runmode to the appname
provided to NewRelic (uses a dash as separator). If not defined, the NewRelic appname is the beego appname.

    NewrelicAppnameRunmode=prepend

reports runmode-appname to NewRelic

    NewrelicAppnameRunmode=append

reports appname-runmode to NewRelic
