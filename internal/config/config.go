package config

import (
  "fmt"
  "github.com/jessevdk/go-flags"
  "go.uber.org/zap"
  "go.uber.org/zap/zapcore"
  "time"
)

type OptionsSrv struct {
  Host   string `short:"h" long:"host" description:"хост" default:"localhost" env:"HOST"`
  Port   string `short:"p" long:"port" description:"порт" default:"3000" env:"PORT"`
  Log    string `long:"logger-create" description:"logger-create output" default:"debug" env:"LOG"`
  DbHost string `long:"dbhost" description:"the db server host" default:"localhost" env:"DB_HOST"`
  DbPort string `long:"dbport" description:"the db server port" default:"5432" env:"DB_PORT"`
  PgUser string `long:"pguser" description:"the db user" default:"user_postgres" env:"POSTGRES_USER"`
  PgPass string `long:"pgpass" description:"the db pass" default:"pass" env:"POSTGRES_PASSWORD"`
  DbName string `long:"dbname" description:"the db name" default:"test" env:"POSTGRES_DB"`
}

type ConfSrv struct {
  Options OptionsSrv
  Logger  *zap.Logger
}

func InitConfServ() (*ConfSrv, error) {
  var conf ConfSrv
  var opts OptionsSrv
  parser := flags.NewParser(&opts, flags.Default)
  _, err := parser.Parse()
  if err != nil {
    return nil, err
  }
  logger := initLogger(opts.Log)
  conf.Options = opts
  conf.Logger = logger
  return &conf, nil
}

func initLogger(option string) *zap.Logger {
  config := zap.NewDevelopmentConfig()
  config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

  if option == "prod" {
    config = zap.NewProductionConfig()
    config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
  }

  config.EncoderConfig.TimeKey = "timestamp"
  config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
    enc.AppendString(t.UTC().Format(time.RFC3339))
  }
  logger, err := config.Build()
  if err != nil {
    panic("cannot initialize logger-create")
  }
  return logger
}

func (o OptionsSrv) DbString(schema string) string {
  return fmt.Sprintf("%s://%s:%s@%s:%s/%s", schema, o.PgUser, o.PgPass, o.DbHost, o.DbPort, o.DbName)
}
