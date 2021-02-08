module github.com/go-vela/worker

replace github.com/go-vela/pkg-queue => github.com/JordanSussman/pkg-queue v0.6.0-rc1.0.20210208212131-07a11f7fab70

go 1.13

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-vela/pkg-executor v0.7.0
	github.com/go-vela/pkg-queue v0.7.0
	github.com/go-vela/pkg-runtime v0.7.0
	github.com/go-vela/sdk-go v0.7.1-0.20210208200135-738ada5b7003
	github.com/go-vela/types v0.7.1-0.20210204153653-939416ae12ed
	github.com/joho/godotenv v1.3.0
	github.com/prometheus/client_golang v1.9.0
	github.com/sirupsen/logrus v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
)
