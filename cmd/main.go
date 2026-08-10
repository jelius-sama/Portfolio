// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package main

import (
    "net/url"
    "os"
    "path/filepath"
    "strings"

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
    var env []types.Env

    var parsedURL, err = url.Parse(AssetCDNHost)
    if err != nil {
        logger.Panic(err)
    }

    if u, err := url.Parse(Host); err != nil {
        logger.Panic(err)
    } else {
        // Grab everything up to the very first dot
        if idx := strings.Index(u.Host, "."); idx != -1 {
            env = append(env,
                types.Env{Key: types.EVDomainname.Get().Key, Value: u.Host[:idx]},
            )
        } else {
            env = append(env,
                types.Env{Key: types.EVDomainname.Get().Key, Value: u.Hostname()},
            )
        }
    }

    var dataDir string = os.Getenv("XDG_DATA_HOME")
    if len(dataDir) == 0 {
        var home, err = os.UserHomeDir()
        if err != nil {
            logger.Panic(err)
        }
        dataDir = filepath.Join(home, ".local", "share", parsedURL.Hostname())
    } else {
        dataDir = filepath.Join(dataDir, parsedURL.Hostname())
    }

    env = append(env,
        types.Env{Key: types.EVEnv.Get().Key, Value: Environment},
        types.Env{Key: types.EVHostname.Get().Key, Value: Host},
        types.Env{Key: types.EVAssetCDNHostname.Get().Key, Value: AssetCDNHost},
        types.Env{Key: types.EVVersion.Get().Key, Value: Version},
        types.Env{Key: types.EVPort.Get().Key, Value: Port},
        types.Env{Key: types.EVDataDir.Get().Key, Value: dataDir},
    )

    types.InitEnv(env)

    logger.Configure(logger.Cnf{
        IsDev: logger.IsDev{
            EnvironmentVariable: new(types.EVEnv.Get().Key),
            ExpectedValue:       new(types.EMDev.String()),
        },
        UseSyslog: os.Getenv(types.EVEnv.Get().Key) != types.EMDev.String(),
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

    logger.Info("Server Data Directory:", dataDir)
    logger.Info("SQLite DB Path:", dbPath)
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

