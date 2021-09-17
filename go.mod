module github.com/go-vela/worker

go 1.15

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-vela/pkg-executor v0.8.1
	github.com/go-vela/pkg-queue v0.9.0
	github.com/go-vela/pkg-runtime v0.9.0
	github.com/go-vela/sdk-go v0.9.0
	github.com/go-vela/types v0.9.0
	github.com/joho/godotenv v1.3.0
	github.com/prometheus/client_golang v1.11.0
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

replace github.com/go-vela/sdk-go => ../sdk-go

replace github.com/go-vela/pkg-executor => ../pkg-executor
