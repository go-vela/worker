module github.com/go-vela/worker

go 1.13

require (
	github.com/coreos/go-semver v0.3.0
	github.com/gin-gonic/gin v1.5.0
	github.com/go-vela/pkg-executor v0.0.0-20200409152007-dd0ea738d48e
	github.com/go-vela/pkg-queue v0.0.0-20200324143217-040845faaf50
	github.com/go-vela/pkg-runtime v0.0.0-20200409170123-8220bae0dff7
	github.com/go-vela/sdk-go v0.3.1-0.20200316181126-22974be2a711
	github.com/go-vela/types v0.3.1-0.20200408124446-5750ec2cac11
	github.com/joho/godotenv v1.3.0
	github.com/prometheus/client_golang v1.2.1
	github.com/sirupsen/logrus v1.5.0
	github.com/urfave/cli/v2 v2.2.0
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637
)

replace github.com/go-vela/pkg-executor => github.com/GregoryDosh/pkg-executor v0.0.0-20200410001857-bf649792150f

replace github.com/go-vela/pkg-queue => github.com/GregoryDosh/pkg-queue v0.0.0-20200409234805-1d6182a8eba4

replace github.com/go-vela/pkg-runtime => github.com/GregoryDosh/pkg-runtime v0.0.0-20200410000527-9a1b10bc1f14
