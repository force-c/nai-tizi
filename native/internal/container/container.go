package container

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/gcc798/quick.admin/internal/config"
	"github.com/gcc798/quick.admin/internal/database"
	"github.com/gcc798/quick.admin/internal/database/migrations"
	"github.com/gcc798/quick.admin/internal/domain/model"
	"github.com/gcc798/quick.admin/internal/logger"
	"github.com/gcc798/quick.admin/internal/messaging/websocket"
	"github.com/gcc798/quick.admin/pkg/captcha"
	"github.com/gcc798/quick.admin/pkg/jwt"
	"github.com/gcc798/quick.admin/pkg/mqtt"
	"github.com/gcc798/quick.admin/pkg/rabbitmq"
	redisclient "github.com/gcc798/quick.admin/pkg/redis"
	"github.com/gcc798/quick.admin/pkg/s3"
	"github.com/gcc798/quick.admin/pkg/scheduler"
	"github.com/gcc798/quick.admin/pkg/storage"
	"github.com/gcc798/quick.admin/pkg/thirdparty/email"
	"github.com/gcc798/quick.admin/pkg/thirdparty/sms"
	"github.com/gcc798/quick.admin/pkg/thirdparty/wechat"
	goredis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Component 统一组件接口
type Component interface {
	Name() string
	Start() error
	Stop() error
}

// Container 依赖注入容器接口
type Container interface {
	GetConfig() *config.Config
	GetViper() *viper.Viper
	GetDB() *gorm.DB
	GetRedis() *goredis.Client
	GetJWT() *jwt.Jwt
	GetLogger() logger.Logger
	GetMQTT() *mqtt.Client
	GetRabbitMQProducer() *rabbitmq.ProducerService
	GetWeChat() *wechat.Manager
	GetSMS() *sms.Manager
	GetEmail() *email.Manager
	GetS3() *s3.Manager
	GetStorageManager() storage.StorageManager
	GetWebSocketHub() *websocket.Hub
	GetScheduler() *scheduler.Scheduler
	GetCaptchaManager() *captcha.CaptchaManager
	RegisterComponent(comp Component)
	Start() error
	Stop()
}

type container struct {
	config         *config.Config
	viper          *viper.Viper
	db             *gorm.DB
	redis          *goredis.Client
	jwt            *jwt.Jwt
	logger         logger.Logger
	mqttClient     *mqtt.Client
	rabbitMQ       *rabbitmq.Manager
	wechatManager  *wechat.Manager
	smsManager     *sms.Manager
	emailManager   *email.Manager
	s3Manager      *s3.Manager
	storageManager storage.StorageManager
	wsHub          *websocket.Hub
	sched          *scheduler.Scheduler
	captchaManager *captcha.CaptchaManager

	components []Component
}

const weChatXcxConfigCode = "WechatXcxCfg"

// NewEmpty 创建一个空容器，调用方可以按需初始化指定组件。
func NewEmpty(cfg *config.Config, v *viper.Viper, log logger.Logger) *container {
	return &container{
		config:     cfg,
		viper:      v,
		logger:     log,
		components: make([]Component, 0),
	}
}

// New 创建新的容器实例
func New(cfg *config.Config, v *viper.Viper, log logger.Logger) (Container, error) {
	c := &container{
		config:     cfg,
		viper:      v,
		logger:     log,
		components: make([]Component, 0),
	}

	// 1. 初始化基础组件
	if err := c.initDB(); err != nil {
		return nil, err
	}
	if err := c.initRedis(); err != nil {
		return nil, err
	}
	c.initJWT()

	// 2. 初始化业务组件
	if err := c.initMQTT(); err != nil {
		return nil, err
	}
	if err := c.initRabbitMQ(); err != nil {
		return nil, err
	}
	c.initThirdParty()
	c.initStorageManager()
	c.initWebSocket()
	c.initCaptchaManager()

	// 3. 初始化调度器
	c.initScheduler()

	return c, nil
}

// InitDBOnly 仅初始化数据库连接，适合一次性工具复用现有连库逻辑。
func (c *container) InitDBOnly() error {
	return c.initDB()
}

// RegisterComponent 执行业务逻辑。
func (c *container) RegisterComponent(comp Component) {
	c.components = append(c.components, comp)
}

// initDB 初始化数据库
func (c *container) initDB() error {
	dsn := c.config.Database.DSN
	slowThreshold := time.Second
	if c.config.Database.SlowThreshold > 0 {
		slowThreshold = time.Duration(c.config.Database.SlowThreshold) * time.Millisecond
	}
	gormLogger := gormlogger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             slowThreshold,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(c.config.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.config.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.config.Database.ConnMaxLifetimeMinutes) * time.Minute)

	// 注册 ID 生成插件
	idGenPlugin := &database.IDGenPlugin{}
	if err := db.Use(idGenPlugin); err != nil {
		c.logger.Warn("failed to register ID generation plugin", zap.Error(err))
	}

	// GORM AutoMigrate 配置
	if c.config.Database.AutoMigrate {
		c.logger.Info("starting database auto migration...")
		if err := db.AutoMigrate(
			&model.User{},
			&model.DictData{},
			&model.LoginLog{},
			&model.OperLog{},
			&model.AuthClient{},
			&model.BuMessageRetry{},
			&model.BuMessageRetryLog{},
			&model.Role{},
			&model.Menu{},
			&model.Org{},
			&model.Config{},
			&model.StorageEnv{},
			&model.Attachment{},
			&model.MUserRole{},
			&model.MRoleMenu{},
			&model.ApiPermission{},
			&model.MRoleApiPermission{},
			&model.MUserApiPermission{},
		); err != nil {
			return fmt.Errorf("failed to auto migrate: %w", err)
		}
		c.logger.Info("database auto migration completed")
	}

	// 当前工程仍通过 AutoMigrate 创建基础表，因此在其后执行 Goose，兼容全新数据库与存量数据库。
	c.logger.Info("starting goose database migrations...")
	if err := migrations.Up(sqlDB); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	c.logger.Info("goose database migrations completed")

	c.db = db
	return nil
}

// initRedis 初始化Redis
func (c *container) initRedis() error {
	redisClient := redisclient.NewRedis(c.config.Redis.Addr, c.config.Redis.Password, c.config.Redis.DB)
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	c.redis = redisClient
	return nil
}

// initJWT 初始化JWT
func (c *container) initJWT() {
	c.jwt = jwt.New(c.config.JWT.Secret, int64(c.config.JWT.Expire))
}

// initMQTT 初始化MQTT
func (c *container) initMQTT() error {
	if !c.config.MQTT.Enabled {
		return nil
	}
	client, err := mqtt.NewClient(&mqtt.Config{
		Broker:   c.config.MQTT.Broker,
		ClientID: c.config.MQTT.ClientID,
		Username: c.config.MQTT.Username,
		Password: c.config.MQTT.Password,
		QoS:      c.config.MQTT.QoS,
	}, c.logger)
	if err != nil {
		return fmt.Errorf("failed to create MQTT client: %w", err)
	}

	c.mqttClient = client
	c.RegisterComponent(client)
	return nil
}

// initRabbitMQ 初始化RabbitMQ
func (c *container) initRabbitMQ() error {
	if !c.config.RabbitMQ.Enabled {
		return nil
	}
	manager, err := rabbitmq.NewManager(&rabbitmq.Config{
		URL:     c.config.RabbitMQ.URL,
		Enabled: c.config.RabbitMQ.Enabled,
	}, c.logger)
	if err != nil {
		c.logger.Warn("failed to create RabbitMQ manager", zap.Error(err))
		return nil // 允许失败，不阻断启动
	}
	c.rabbitMQ = manager
	c.RegisterComponent(manager)
	return nil
}

// initThirdParty 初始化第三方服务
func (c *container) initThirdParty() {
	c.loadWeChatXcxConfig()
	if c.config.WeChat.Enabled {
		c.wechatManager = wechat.NewManager(wechat.Config{
			Enabled: c.config.WeChat.Enabled,
			AppID:   c.config.WeChat.AppID,
			Secret:  c.config.WeChat.Secret,
		}, c.logger, c.redis)
	}

	if c.config.SMS.Enabled {
		smsManager, err := sms.NewManager(sms.Config{
			AccessKeyId:     c.config.SMS.AccessKeyId,
			AccessKeySecret: c.config.SMS.AccessKeySecret,
			SignName:        c.config.SMS.SignName,
			TemplateCode:    c.config.SMS.TemplateCode,
		}, c.redis, c.logger)
		if err != nil {
			c.logger.Warn("failed to create SMS service", zap.Error(err))
		} else {
			c.smsManager = smsManager
		}
	}

	if c.config.Email.Enabled {
		emailManager, err := email.NewManager(email.Config{
			Host:     c.config.Email.Host,
			Port:     c.config.Email.Port,
			Username: c.config.Email.Username,
			Password: c.config.Email.Password,
			From:     c.config.Email.From,
		}, c.redis, c.logger)
		if err != nil {
			c.logger.Warn("failed to create email service", zap.Error(err))
		} else {
			c.emailManager = emailManager
		}
	}

	if c.config.S3.Enabled {
		s3Manager, err := s3.NewManager(&s3.Config{
			Enabled:         c.config.S3.Enabled,
			Endpoint:        c.config.S3.Endpoint,
			AccessKeyID:     c.config.S3.AccessKeyID,
			SecretAccessKey: c.config.S3.SecretAccessKey,
			Region:          c.config.S3.Region,
			Bucket:          c.config.S3.Bucket,
			UseSSL:          c.config.S3.UseSSL,
			ForcePathStyle:  c.config.S3.ForcePathStyle,
		}, c.logger)
		if err != nil {
			c.logger.Warn("failed to create S3 manager", zap.Error(err))
		} else {
			c.s3Manager = s3Manager
		}
	}
}

func (c *container) loadWeChatXcxConfig() {
	var item model.Config
	err := c.db.Where("code = ?", weChatXcxConfigCode).First(&item).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.logger.Warn("failed to load wechat xcx config from database", zap.Error(err))
		}
		return
	}

	dbConfig := c.config.WeChat
	if err := json.Unmarshal(item.Data, &dbConfig); err != nil {
		c.logger.Warn("failed to parse wechat xcx config", zap.Error(err))
		return
	}
	c.config.WeChat = dbConfig
}

// initStorageManager 初始化存储管理器
func (c *container) initStorageManager() {
	// 创建存储管理器
	c.storageManager = storage.NewStorageManager(c.db, c.logger)

	// 注册存储类型工厂
	c.storageManager.RegisterStorageType("s3", storage.NewS3StorageFactory())
	c.storageManager.RegisterStorageType("local", storage.NewLocalStorageFactory())

	c.logger.Info("storage manager initialized successfully")
}

// initCaptchaManager 初始化验证码管理器
func (c *container) initCaptchaManager() {
	manager := captcha.NewCaptchaManager()

	// 注册图形验证码提供者
	if c.config.Captcha.Image.Enabled {
		imageConfig := &captcha.ImageCaptchaConfig{
			Enabled: c.config.Captcha.Image.Enabled,
			Length:  c.config.Captcha.Image.Length,
			Width:   c.config.Captcha.Image.Width,
			Height:  c.config.Captcha.Image.Height,
			Expire:  c.config.Captcha.Image.Expire,
		}
		imageProvider := captcha.NewImageCaptchaProvider(imageConfig, c.redis)
		manager.RegisterProvider(imageProvider)
	}

	// 注册短信验证码提供者
	if c.config.Captcha.SMS.Enabled {
		if c.smsManager == nil {
			c.logger.Warn("SMS captcha enabled but SMS service not configured")
		} else {
			smsConfig := &captcha.SMSCaptchaConfig{
				Enabled:  c.config.Captcha.SMS.Enabled,
				Length:   c.config.Captcha.SMS.Length,
				Expire:   c.config.Captcha.SMS.Expire,
				Template: c.config.Captcha.SMS.Template,
				Provider: c.config.Captcha.SMS.Provider,
			}
			smsAdapter := captcha.NewSMSManagerAdapter(c.smsManager)
			smsProvider := captcha.NewSMSCaptchaProvider(smsConfig, c.redis, smsAdapter)
			manager.RegisterProvider(smsProvider)
		}
	}

	// 注册邮箱验证码提供者
	if c.config.Captcha.Email.Enabled {
		if c.emailManager == nil {
			c.logger.Warn("Email captcha enabled but email service not configured")
		} else {
			emailConfig := &captcha.EmailCaptchaConfig{
				Enabled:  c.config.Captcha.Email.Enabled,
				Length:   c.config.Captcha.Email.Length,
				Expire:   c.config.Captcha.Email.Expire,
				Template: c.config.Captcha.Email.Template,
			}
			emailAdapter := captcha.NewEmailManagerAdapter(c.emailManager)
			emailProvider := captcha.NewEmailCaptchaProvider(emailConfig, c.redis, emailAdapter)
			manager.RegisterProvider(emailProvider)
		}
	}

	c.captchaManager = manager
	c.logger.Info("captcha manager initialized successfully")
}

// initWebSocket 初始化WebSocket
func (c *container) initWebSocket() {
	if !c.config.WebSocket.Enabled {
		return
	}
	c.wsHub = websocket.NewHub(c.logger)
	c.RegisterComponent(c.wsHub)
}

// initScheduler 初始化调度器
func (c *container) initScheduler() {
	if !c.config.Scheduler.Enabled {
		return
	}
	c.sched = scheduler.New(c.logger)
	c.RegisterComponent(c.sched)
}

// GetConfig 获取业务数据。
func (c *container) GetConfig() *config.Config {
	return c.config
}

// GetViper 获取业务数据。
func (c *container) GetViper() *viper.Viper { return c.viper }

// GetDB 获取业务数据。
func (c *container) GetDB() *gorm.DB {
	return c.db
}

// GetRedis 获取业务数据。
func (c *container) GetRedis() *goredis.Client {
	return c.redis
}

// GetJWT 获取业务数据。
func (c *container) GetJWT() *jwt.Jwt {
	return c.jwt
}

// GetLogger 获取业务数据。
func (c *container) GetLogger() logger.Logger {
	return c.logger
}

// GetMQTT 获取业务数据。
func (c *container) GetMQTT() *mqtt.Client {
	return c.mqttClient
}

// GetRabbitMQProducer 获取业务数据。
func (c *container) GetRabbitMQProducer() *rabbitmq.ProducerService {
	if c.rabbitMQ == nil {
		return nil
	}
	return c.rabbitMQ.GetProducer()
}

// GetWeChat 获取业务数据。
func (c *container) GetWeChat() *wechat.Manager {
	return c.wechatManager
}

// GetSMS 获取业务数据。
func (c *container) GetSMS() *sms.Manager {
	return c.smsManager
}

// GetEmail 获取业务数据。
func (c *container) GetEmail() *email.Manager {
	return c.emailManager
}

// GetS3 获取业务数据。
func (c *container) GetS3() *s3.Manager {
	return c.s3Manager
}

// GetStorageManager 获取业务数据。
func (c *container) GetStorageManager() storage.StorageManager {
	return c.storageManager
}

// GetWebSocketHub 获取业务数据。
func (c *container) GetWebSocketHub() *websocket.Hub {
	return c.wsHub
}

// GetScheduler 获取业务数据。
func (c *container) GetScheduler() *scheduler.Scheduler {
	return c.sched
}

// GetCaptchaManager 获取业务数据。
func (c *container) GetCaptchaManager() *captcha.CaptchaManager {
	return c.captchaManager
}

// Start 启动组件。
func (c *container) Start() error {
	for _, comp := range c.components {
		c.logger.Info("starting component", zap.String("name", comp.Name()))
		if err := comp.Start(); err != nil {
			c.logger.Error("failed to start component", zap.String("name", comp.Name()), zap.Error(err))
			return err
		}
	}
	c.logger.Info("all components started successfully")
	return nil
}

// Stop 停止组件。
func (c *container) Stop() {
	// 反向停止
	for i := len(c.components) - 1; i >= 0; i-- {
		comp := c.components[i]
		c.logger.Info("stopping component", zap.String("name", comp.Name()))
		if err := comp.Stop(); err != nil {
			c.logger.Error("failed to stop component", zap.String("name", comp.Name()), zap.Error(err))
		}
	}
	c.logger.Info("all components stopped")
}
