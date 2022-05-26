package data

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/extra/redisotel"
    "github.com/go-redis/redis/v8"
    "github.com/google/wire"
    "github.com/jassue/gin-wire/app/service"
    "github.com/jassue/gin-wire/config"
    "github.com/jassue/gin-wire/utils"
    "github.com/sony/sonyflake"
    "go.uber.org/zap"
    "gopkg.in/natefinch/lumberjack.v2"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "io"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
    NewData, NewDB, NewRedis, NewTransaction,
    NewUserRepo, NewJwtRepo, NewMediaRepo,
    )

// Data .
type Data struct {
    db *gorm.DB
    rdb *redis.Client
    sf *sonyflake.Sonyflake
}

// NewData .
func NewData(logger *zap.Logger, db *gorm.DB, rdb *redis.Client, sf *sonyflake.Sonyflake) (*Data, func(), error) {
    cleanup := func() {
        logger.Info("closing the data resources")
    }

    return &Data{db: db, rdb: rdb, sf: sf}, cleanup, nil
}

// NewDB .
func NewDB(conf *config.Configuration, gLog *zap.Logger) *gorm.DB {
    if conf.Database.Driver != "mysql" {
        panic(conf.Database.Driver + " driver is not supported")
    }

    var writer io.Writer
    var logMode logger.LogLevel

    // 是否启用日志文件
    if conf.Database.EnableFileLogWriter {
        // 自定义 Writer
        writer = &lumberjack.Logger{
            Filename:   filepath.Join(utils.RootPath(), conf.Log.RootDir, conf.Database.LogFilename),
            MaxSize:    conf.Log.MaxSize,
            MaxBackups: conf.Log.MaxBackups,
            MaxAge:     conf.Log.MaxAge,
            Compress:   conf.Log.Compress,
        }
    } else {
        // 默认 Writer
        writer = os.Stdout
    }

    switch conf.Database.LogMode {
    case "silent":
        logMode = logger.Silent
    case "error":
        logMode = logger.Error
    case "warn":
        logMode = logger.Warn
    case "info":
        logMode = logger.Info
    default:
        logMode = logger.Info
    }

    newLogger := logger.New(
        log.New(writer, "\r\n", log.LstdFlags), // io writer
        logger.Config{
            SlowThreshold: time.Second, // 慢查询 SQL 阈值
            Colorful:      !conf.Database.EnableFileLogWriter, // 禁用彩色打印
            IgnoreRecordNotFoundError: false, // 忽略ErrRecordNotFound（记录未找到）错误
            LogLevel: logMode, // Log lever
        },
    )

    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
        conf.Database.UserName,
        conf.Database.Password,
        conf.Database.Host,
        strconv.Itoa(conf.Database.Port),
        conf.Database.Database,
        conf.Database.Charset,
        )
    mysqlConfig := mysql.Config{
        DSN:                       dsn,   // DSN data source name
        DefaultStringSize:         191,   // string 类型字段的默认长度
        DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
        DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
        DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
        SkipInitializeWithVersion: false, // 根据版本自动配置
    }
    if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
        Logger: newLogger, // 使用自定义 Logger
    }); err != nil {
        gLog.Error("failed opening connection to err:", zap.Any("err", err))
        panic("failed to connect database")
    } else {
        sqlDB, _ := db.DB()
        sqlDB.SetMaxIdleConns(conf.Database.MaxIdleConns)
        sqlDB.SetMaxOpenConns(conf.Database.MaxOpenConns)
        return db
    }
}

// NewRedis .
func NewRedis(c *config.Configuration, gLog *zap.Logger) *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     c.Redis.Host + ":" + c.Redis.Port,
        Password: c.Redis.Password, // no password set
        DB:       c.Redis.DB,       // use default DB
    })

    client.AddHook(redisotel.TracingHook{})
    if err := client.Close(); err != nil {
        gLog.Error("redis connect failed, err:", zap.Any("err", err))
    }

    return client
}

type contextTxKey struct{}

func (d *Data) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
   return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
       ctx = context.WithValue(ctx, contextTxKey{}, tx)
       return fn(ctx)
   })
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
   tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
   if ok {
       return tx
   }
   return d.db
}

// NewTransaction .
func NewTransaction(d *Data) service.Transaction {
   return d
}

