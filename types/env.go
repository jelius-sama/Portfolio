package types

import "os"

type Env struct {
    Key   string
    Value string
}

type EnvMode uint8

const (
    EMDev EnvMode = iota
    EMProd
)

func (em EnvMode) String() string {
    switch em {
    case EMDev:
        return "development"
    case EMProd:
        return "production"
    default:
        return "unknown"
    }
}

type EnvVal uint8

const (
    EVEnv EnvVal = iota
    EVHostname
    EVAssetCDNHostname
    EVVersion
    EVPort
    EVDataDir
)

func (ek EnvVal) Get() Env {
    switch ek {
    case EVEnv:
        return Env{Key: "ENV", Value: os.Getenv("ENV")}
    case EVHostname:
        return Env{Key: "HOSTNAME", Value: os.Getenv("HOSTNAME")}
    case EVAssetCDNHostname:
        return Env{Key: "ASSET_CDN_HOSTNAME", Value: os.Getenv("ASSET_CDN_HOSTNAME")}
    case EVVersion:
        return Env{Key: "VERSION", Value: os.Getenv("VERSION")}
    case EVPort:
        return Env{Key: "PORT", Value: os.Getenv("PORT")}
    case EVDataDir:
        return Env{Key: "DATA_DIR", Value: os.Getenv("DATA_DIR")}
    default:
        return Env{Key: "", Value: ""}
    }
}

func InitEnv[V []Env | Env](x V) {
    switch env := any(x).(type) {
    case Env:
        if _, present := os.LookupEnv(env.Key); !present {
            os.Setenv(env.Key, env.Value)
        }
    case []Env:
        for i := range env {
            if _, present := os.LookupEnv(env[i].Key); !present {
                os.Setenv(env[i].Key, env[i].Value)
            }
        }
    }
}

