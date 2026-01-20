package mqtt

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	logging "github.com/force-c/nai-tizi/internal/logger"
	"go.uber.org/zap"
)

type Subscription struct {
	Topic   string
	Handler mqtt.MessageHandler
}

type Client struct {
	client           mqtt.Client
	logger           logging.Logger
	broker, clientID string
	qos              byte
	subscriptions    []Subscription
}

type Config struct {
	Broker, ClientID, Username, Password string
	QoS                                  byte
}

func NewClient(cfg *Config, logger logging.Logger) (*Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientID)
	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) { logger.Warn("MQTT connection lost", zap.Error(err)) })
	opts.SetOnConnectHandler(func(client mqtt.Client) { logger.Info("MQTT connected successfully") })
	c := &Client{client: mqtt.NewClient(opts), logger: logger, broker: cfg.Broker, clientID: cfg.ClientID, qos: cfg.QoS}
	return c, nil
}

func (c *Client) AddSubscription(topic string, handler mqtt.MessageHandler) {
	c.subscriptions = append(c.subscriptions, Subscription{Topic: topic, Handler: handler})
}

func (c *Client) Name() string {
	return "MQTT Client"
}

func (c *Client) Start() error {
	if err := c.Connect(); err != nil {
		return err
	}
	for _, sub := range c.subscriptions {
		if err := c.Subscribe(sub.Topic, sub.Handler); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Stop() error {
	c.Disconnect()
	return nil
}

func (c *Client) Connect() error {
	token := c.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}
	c.logger.Info("MQTT client connected successfully")
	return nil
}
func (c *Client) Disconnect() {
	c.logger.Info("disconnecting from MQTT broker")
	c.client.Disconnect(250)
}
func (c *Client) IsConnected() bool { return c.client.IsConnected() }
func (c *Client) Subscribe(topic string, callback mqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, c.qos, callback)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}
	return nil
}
func (c *Client) Unsubscribe(topic string) error {
	token := c.client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe from topic %s: %w", topic, token.Error())
	}
	return nil
}
func (c *Client) Publish(topic string, payload string) error {
	if !c.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}
	token := c.client.Publish(topic, c.qos, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish message: %w", token.Error())
	}
	return nil
}
func (c *Client) PublishControl(netType, mac, sn string, payload string) error {
	topic := fmt.Sprintf("NTZ/%s/%s/%s", netType, mac, sn)
	return c.Publish(topic, payload)
}
