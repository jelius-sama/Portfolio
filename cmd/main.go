package main

import (
    "net/url"
    "os"
    "path/filepath"

    "git.jelius.dev/jelius-sama/Portfolio/db"
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
    var parsedURL, err = url.Parse(types.EVAssetCDNHostname.Get().Value)
    if err != nil {
        logger.Panic(err)
    }

    var dataDir string = os.Getenv("XDG_DATA_HOME")
    if len(dataDir) == 0 {
        var home, err = os.UserHomeDir()
        if err != nil {
            logger.Panic(err)
        }
        dataDir = filepath.Join(home, ".local", "share", parsedURL.Hostname())
    }

    var env []types.Env = []types.Env{
        {Key: types.EVEnv.Get().Key, Value: Environment},
        {Key: types.EVHostname.Get().Key, Value: Host},
        {Key: types.EVAssetCDNHostname.Get().Key, Value: AssetCDNHost},
        {Key: types.EVVersion.Get().Key, Value: Version},
        {Key: types.EVPort.Get().Key, Value: Port},
        {Key: types.EVDataDir.Get().Key, Value: dataDir},
    }

    types.InitEnv(env)

    logger.Configure(logger.Cnf{
        IsDev: logger.IsDev{
            EnvironmentVariable: new(types.EVEnv.Get().Key),
            ExpectedValue:       new(types.EMDev.String()),
        },
        UseSyslog: false,
        // UseSyslog: os.Getenv(types.EVEnv.Get().Key) != types.EMDev.String(),
    })

    var dbPath string
    if types.EVEnv.Get().Value == types.EMDev.String() {
        cwd, err := os.Getwd()
        if err != nil {
            logger.Fatal("Failed to determine working directory:", err.Error())
        }
        dbPath = filepath.Join(cwd, "db.sqlite3")
    } else {
        dbPath = filepath.Join(dataDir, "db.sqlite3")
    }

    logger.Info("Using DB Path:", dbPath)
    err = db.InitDB(dbPath)
    if err != nil {
        logger.Fatal("Failed to initialize database:", err.Error())
    }
}

func main() {
    defer db.CloseDB() // if we are in main we can assume that `init()` was successful.

    var cnf fiber.Config = fiber.Config{
        ErrorHandler: middleware.ErrHandler,
    }

    var app *fiber.App = fiber.New(cnf)
    Router(app)

    logger.Fatal(app.Listen(types.EVPort.Get().Value))
}

