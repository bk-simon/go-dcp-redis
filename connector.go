package dcpredis

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"

	jsoniter "github.com/json-iterator/go"

	dcpCouchbase "github.com/Trendyol/go-dcp/couchbase"

	"github.com/Trendyol/go-dcp-redis/metric"

	"github.com/Trendyol/go-dcp"
	"github.com/Trendyol/go-dcp-redis/config"
	"github.com/Trendyol/go-dcp-redis/couchbase"
	"github.com/Trendyol/go-dcp-redis/redis/bulk"
	"github.com/Trendyol/go-dcp/logger"
	"github.com/Trendyol/go-dcp/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Connector interface {
	Start()
	Close()
	GetDcpClient() dcpCouchbase.Client
}

type connector struct {
	dcp    dcp.Dcp
	mapper Mapper
	config *config.Connector
	bulk   *bulk.Bulk
}

func (c *connector) Start() {
	go func() {
		<-c.dcp.WaitUntilReady()
		c.bulk.StartBulk()
	}()
	c.dcp.Start()
}

func (c *connector) Close() {
	c.dcp.Close()
	c.bulk.Close()
}

func (c *connector) GetDcpClient() dcpCouchbase.Client {
	return c.dcp.GetClient()
}

func (c *connector) listener(ctx *models.ListenerContext) {
	// Initialize ListenerTrace for current listen operation
	listenerTrace := ctx.ListenerTracerComponent.InitializeListenerTrace("Listen", nil)
	defer listenerTrace.Finish()

	// Get redis client
	redisClient := c.bulk.GetRedisClient()

	var e couchbase.Context
	switch event := ctx.Event.(type) {
	case models.DcpMutation:
		params := couchbase.EventParams{
			Key:            event.Key,
			Value:          event.Value,
			CollectionName: event.CollectionName,
			EventTime:      event.EventTime,
			Cas:            event.Cas,
			VbID:           event.VbID,
		}
		e = couchbase.NewMutateEventContext(listenerTrace, redisClient, params)
	case models.DcpExpiration:
		params := couchbase.EventParams{
			Key:            event.Key,
			Value:          nil,
			CollectionName: event.CollectionName,
			EventTime:      event.EventTime,
			Cas:            event.Cas,
			VbID:           event.VbID,
		}
		e = couchbase.NewExpireEventContext(listenerTrace, redisClient, params)
	case models.DcpDeletion:
		params := couchbase.EventParams{
			Key:            event.Key,
			Value:          nil,
			CollectionName: event.CollectionName,
			EventTime:      event.EventTime,
			Cas:            event.Cas,
			VbID:           event.VbID,
		}
		e = couchbase.NewDeleteEventContext(listenerTrace, redisClient, params)
	default:
		return
	}

	actions := c.mapper(e)

	if len(actions) == 0 {
		ctx.Ack()
		return
	}

	c.bulk.AddActions(ctx, e.Event.EventTime, actions)
}

type ConnectorBuilder struct {
	mapper Mapper
	config any
}

func newConnectorConfigFromPath(path string) (*config.Connector, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c config.Connector
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func newConfig(cf any) (*config.Connector, error) {
	switch v := cf.(type) {
	case *config.Connector:
		return v, nil
	case config.Connector:
		return &v, nil
	case string:
		return newConnectorConfigFromPath(v)
	default:
		return nil, errors.New("invalid config")
	}
}

func newConnector(cf any, mapper Mapper) (Connector, error) {
	cfg, err := newConfig(cf)
	if err != nil {
		return nil, err
	}
	cfg.ApplyDefaults()

	connector := &connector{
		mapper: mapper,
		config: cfg,
	}

	dcp, err := dcp.NewDcp(&cfg.Dcp, connector.listener)
	if err != nil {
		logger.Log.Error("Dcp error: %v", err)
		return nil, err
	}

	copyOfConfig := cfg.Redis
	printConfiguration(copyOfConfig)

	dcpConfig := dcp.GetConfig()
	dcpConfig.Checkpoint.Type = "manual"

	connector.dcp = dcp

	connector.bulk, err = bulk.NewBulk(cfg, dcp.Commit)
	if err != nil {
		return nil, err
	}

	connector.dcp.SetEventHandler(
		&DcpEventHandler{
			isFinite: dcpConfig.IsDcpModeFinite(),
			bulk:     connector.bulk,
		})

	metricCollector := metric.NewMetricCollector(connector.bulk)
	dcp.SetMetricCollectors(metricCollector)

	SetCollectionKeyMappings(&cfg.Redis.CollectionKeyMapping)

	return connector, nil
}

func NewConnectorBuilder(config any) ConnectorBuilder {
	return ConnectorBuilder{
		config: config,
		mapper: DefaultMapper,
	}
}

func (c ConnectorBuilder) SetMapper(mapper Mapper) ConnectorBuilder {
	c.mapper = mapper
	return c
}

func (c ConnectorBuilder) Build() (Connector, error) {
	return newConnector(c.config, c.mapper)
}

func (c ConnectorBuilder) SetLogger(logrus *logrus.Logger) ConnectorBuilder {
	logger.Log = &logger.Loggers{
		Logrus: logrus,
	}
	return c
}

func printConfiguration(config config.Redis) {
	config.Password = "*****"
	configJSON, _ := jsoniter.Marshal(config)

	dst := &bytes.Buffer{}
	if err := json.Compact(dst, configJSON); err != nil {
		logger.Log.Error("error while print redis configuration, err: %v", err)
		panic(err)
	}

	logger.Log.Info("using redis config: %v", dst.String())
}
