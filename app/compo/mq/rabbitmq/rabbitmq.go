package rabbitmq

import (
    "context"
    "errors"
    "fmt"
    "github.com/jassue/gin-wire/app/compo/mq"
    "github.com/jassue/gin-wire/app/compo/mq/event"
    "github.com/jassue/gin-wire/config"
    amqp "github.com/rabbitmq/amqp091-go"
    "go.uber.org/zap"
    "strconv"
    "sync"
    "time"
)

type connManager struct {
    connection *amqp.Connection
    conf config.RabbitMQ
    m sync.Mutex
}

func NewConnManager(c *config.Configuration, log *zap.Logger) *connManager {
    connection, err := amqp.Dial(c.Queue.RabbitMQ.Uri)
    if err != nil {
        log.Error("Failed to connect rabbitmq: " + err.Error())
        panic("Failed to connect rabbitmq: " + err.Error())
    }

    return &connManager{
        connection: connection,
        conf: c.Queue.RabbitMQ,
        m: sync.Mutex{},
    }
}

func (m *connManager) Reconnect() error {
    m.m.Lock()
    if m.connection.IsClosed() {
        connection, err := amqp.Dial(m.conf.Uri)
        if err != nil {
            m.m.Unlock()
            return errors.New("Failed to connect rabbitmq:" + err.Error())
        }
        m.connection = connection
    }
    m.m.Unlock()
    return nil
}

func (m *connManager) NewChannel() (*amqp.Channel, error) {
    channel, err := m.connection.Channel()
    if err != nil {
        return nil, errors.New("Failed to getting channel: " + err.Error())
    }
    return channel, nil
}

type Sender struct {
    channel *amqp.Channel
    connManager *connManager
}

func NewRabbitmqSender(m *connManager) event.Sender {
    channel, err := m.NewChannel()
    if err != nil {
        panic(err)
    }

    // 声明交换机
    if err := channel.ExchangeDeclare(
        m.conf.ExchangeName,
        "direct",
        true,
        false,
        false,
        false,
        nil,
    ); err != nil {
        panic("Exchange declare: " + err.Error())
    }

    // 声明延时交换机
    if err := channel.ExchangeDeclare(
        m.conf.DelayExchangeName,
        "x-delayed-message",
        true,
        false,
        false,
        false,
        amqp.Table{
            "x-delayed-type": "direct",
        },
    ); err != nil {
        panic("Exchange declare: " + err.Error())
    }

    return &Sender{
        channel: channel,
        connManager: m,
    }
}

func (r *Sender) healthCheck() error {
    if r.channel.IsClosed() {
        if err := r.connManager.Reconnect(); err != nil {
            return err
        }
        channel, err := r.connManager.NewChannel()
        if err != nil {
            return err
        }
        r.channel = channel
    }
    return nil
}

func (r *Sender) Send(ctx context.Context, msg event.Event) error {
    // 健康检查
    if err := r.healthCheck(); err != nil {
        return err
    }

    exchangeName := r.connManager.conf.ExchangeName
    // 消息体
    publishing := amqp.Publishing{
        ContentType:     "application/json",
        Body:            msg.Value(),
        Timestamp:       time.Now(),
        DeliveryMode:    amqp.Persistent, // 消息持久化
    }
    if msg.DelaySeconds() > 0 {
        publishing.Headers = amqp.Table{
            "x-delay": msg.DelaySeconds()*1000,
        }
        exchangeName = r.connManager.conf.DelayExchangeName
    }

    // 声明队列
    if _, err := r.channel.QueueDeclare(
        msg.Queue(), // name
        true,  // durable
        false, // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    ); err != nil {
        return errors.New("Failed to declare a queue: " + err.Error())
    }

    // 队列绑定交换机
    if err := r.channel.QueueBind(msg.Queue(), msg.Queue(), exchangeName, false, nil); err != nil {
        return errors.New("Queue Bind: " + err.Error())
    }

    // 发送消息
    err := r.channel.PublishWithContext(
        ctx,
        exchangeName,
        msg.Queue(),
        false,
        false,
        publishing,
        )

    return err
}

func (r *Sender) Close() error {
    return r.channel.Close()
}

type Receiver struct {
    connManager *connManager
    logger *mq.QueueLogger
}

func NewRabbitmqReceiver(m *connManager, logger *mq.QueueLogger) event.Receiver {
    return &Receiver{
        connManager: m,
        logger: logger,
    }
}

func (r *Receiver) Receive(ctx context.Context, e event.Event, handler event.Handler, workerNum int64) error {
    for i := int64(1); i <= workerNum; i++ {
        go func(i string) {
            workerName := e.Queue()+"_worker_"+i
            defer func() {
                if err := recover(); err != nil {
                    r.logger.Error(workerName + " Panic", zap.Any("err", err))
                }
            }()

            isInit := true
            for {
                // 断开重连
                if !isInit {
                    time.Sleep(60*time.Second)
                    if err := r.connManager.Reconnect(); err != nil {
                        r.logger.Error(err.Error())
                        continue
                    }
                } else {
                    isInit = false
                }

                // 创建频道
                channel, err := r.connManager.NewChannel()
                if err != nil {
                    r.logger.Error(err.Error())
                    continue
                }

                // 声明队列
                if _, err := channel.QueueDeclare(
                    e.Queue(), // name
                    true,  // durable
                    false, // delete when unused
                    false, // exclusive
                    false, // no-wait
                    nil,   // arguments
                ); err != nil {
                    r.logger.Error("Failed to declare a queue: " + err.Error())
                    continue
                }

                // 限制消费者获取未确认消息的数量
                if err := channel.Qos(2, 0, false); err != nil {
                    r.logger.Error("Channel Qos: " + err.Error())
                    continue
                }

                deliveries, err := channel.Consume(
                    e.Queue(), // 队列名称
                    workerName, // 消费者名称
                    false, // 是否自动确认消费（false需手动ack）
                    false, // 排他
                    false, // 是否仅接收同一连接的生产者的消息
                    false, // 是否等待服务器确认
                    nil,
                )
                if err != nil {
                    r.logger.Error("Queue Consume: " + err.Error())
                    continue
                }

                // 断开连接信号
                mqErr := channel.NotifyClose(make(chan *amqp.Error))

                r.logger.Info(workerName+" is started")

                isClose := false
                for !isClose {
                    select {
                    case <-ctx.Done():
                        r.logger.Info(workerName + " is stopped")
                        return
                    case <-mqErr:
                        r.logger.Error(workerName + " disconnect")
                        isClose = true
                    case d := <-deliveries:
                        r.logger.Info(workerName + " Processing")

                        if err := handler(ctx, e.Unmarshal(d.Body)); err != nil {
                            r.logger.Error(fmt.Sprintf("%s Processing failed: %v", workerName, err))
                        } else {
                            // 手动ack确认，v.Ack的参数 multiple 表示 是否将当前channel中所有未确认的deliveries都一起确认
                            if err := d.Ack(false); err != nil {
                                r.logger.Error(fmt.Sprintf("%s Deliver ack: %v", workerName, err))
                            } else {
                                r.logger.Info(workerName + " Processed")
                            }
                        }
                        //time.Sleep(300*time.Millisecond)
                    }
                }
            }
        }(strconv.FormatInt(i, 10))
    }

    return nil
}

func (r *Receiver) Close() error {
    return r.connManager.connection.Close()
}
