package main

import (
    "os"

    "git.jelius.dev/jelius-sama/Portfolio/middleware"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

var (
    Environment  string
    Host         string
    AssetCDNHost string
    Version      string
    Port         string
)

func init() {
    var env []types.Env = []types.Env{
        {Key: types.EVEnv.Get().Key, Value: Environment},
        {Key: types.EVHostname.Get().Key, Value: Host},
        {Key: types.EVAssetCDNHostname.Get().Key, Value: AssetCDNHost},
        {Key: types.EVVersion.Get().Key, Value: Version},
        {Key: types.EVPort.Get().Key, Value: Port},
    }

    types.InitEnv(env)

    logger.Configure(logger.Cnf{
        IsDev: logger.IsDev{
            EnvironmentVariable: new(types.EVEnv.Get().Key),
            ExpectedValue:       new(types.EMDev.String()),
        },
        UseSyslog: os.Getenv(types.EVEnv.Get().Key) != types.EMDev.String(),
    })
}

func main() {
    var cnf fiber.Config = fiber.Config{
        ErrorHandler: middleware.ErrHandler,
    }

    var app *fiber.App = fiber.New(cnf)
    Router(app)

    logger.Fatal(app.Listen(types.EVPort.Get().Value))
}

