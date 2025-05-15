package rabbitmq

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"maps"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	_reconnectTimeout      = 5 * time.Second
	_defaultMaxRetries     = 3
	_defaultReconnectDelay = 1 * time.Second
	_maxReconnectBackoff   = 30 * time.Second
	_backoffMultiplier     = 2
)

type Client struct {
	conn                  *amqp.Connection
	queues                map[string]amqp.Queue
	consumers             map[string][]*consumer
	producers             map[string][]*amqp.Channel
	rrIdx                 sync.Map // map[string]*uint64, round-robin index per queue.
	mu                    sync.Mutex
	logger                *zap.Logger
	dsn                   string
	tlsConfig             *tls.Config
	maxRetries            int
	maxReconnectAttempts  int
	initialReconnectDelay time.Duration
	ctx                   context.Context
	cancelFunc            context.CancelFunc
	reconnectMu           sync.Mutex
	isReconnecting        atomic.Bool
}

type consumer struct {
	ch          *amqp.Channel
	doneCh      chan struct{}
	queueName   string
	processFunc func(context.Context, *Message) error
}

// New creates a new RabbitMQ client.
func New(ctx context.Context, dsn string, tlsConf *tls.Config, logger *zap.Logger) (*Client, error) {
	return NewWithConfig(ctx, dsn, tlsConf, logger, _defaultMaxRetries, 0, _defaultReconnectDelay)
}

// NewWithConfig creates a new RabbitMQ client with custom config.
func NewWithConfig(
	ctx context.Context,
	dsn string,
	tlsConf *tls.Config,
	logger *zap.Logger,
	maxRetries int,
	maxReconnectAttempts int,
	initialReconnectDelay time.Duration,
) (*Client, error) {
	if maxRetries < 0 {
		maxRetries = _defaultMaxRetries
	}
	if initialReconnectDelay <= 0 {
		initialReconnectDelay = _defaultReconnectDelay
	}

	conn, err := dial(dsn, tlsConf)
	if err != nil {
		return nil, err
	}

	clientCtx, cancel := context.WithCancel(ctx)
	s := &Client{
		conn:                  conn,
		queues:                make(map[string]amqp.Queue),
		consumers:             make(map[string][]*consumer),
		producers:             make(map[string][]*amqp.Channel),
		logger:                logger,
		dsn:                   dsn,
		tlsConfig:             tlsConf,
		maxRetries:            maxRetries,
		maxReconnectAttempts:  maxReconnectAttempts,
		initialReconnectDelay: initialReconnectDelay,
		ctx:                   clientCtx,
		cancelFunc:            cancel,
	}

	go s.watchConnection()
	go func() {
		<-ctx.Done()
		_ = s.Close()
	}()

	return s, nil
}

func dial(dsn string, tlsConf *tls.Config) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	if tlsConf == nil {
		conn, err = amqp.Dial(dsn)
		if err != nil {
			return nil, errors.Wrap(err, "dial rabbitmq")
		}
	} else {
		conn, err = amqp.DialTLS(dsn, tlsConf)
		if err != nil {
			return nil, errors.Wrap(err, "dial tls rabbitmq")
		}
	}

	return conn, nil
}

func (s *Client) watchConnection() {
	closeChan := s.conn.NotifyClose(make(chan *amqp.Error, 1))
	for {
		select {
		case <-s.ctx.Done():
			return
		case closeErr, ok := <-closeChan:
			if !ok {
				return
			}
			s.logger.Warn("rabbitmq connection closed", zap.Error(closeErr))
			s.reconnect()
			if s.conn != nil && !s.conn.IsClosed() {
				closeChan = s.conn.NotifyClose(make(chan *amqp.Error, 1))
			} else {
				return
			}
		}
	}
}

func (s *Client) reconnect() {
	if !s.isReconnecting.CompareAndSwap(false, true) {
		return
	}
	defer s.isReconnecting.Store(false)

	s.reconnectMu.Lock()
	defer s.reconnectMu.Unlock()

	delay := s.initialReconnectDelay
	attempts := 0

	for {
		if s.maxReconnectAttempts > 0 && attempts >= s.maxReconnectAttempts {
			s.logger.Error("max reconnect attempts reached",
				zap.Int("attempts", attempts),
			)
			return
		}

		select {
		case <-s.ctx.Done():
			return
		case <-time.After(delay):
		}

		attempts++
		s.logger.Info("attempting to reconnect to rabbitmq",
			zap.Int("attempt", attempts),
			zap.Duration("delay", delay),
		)

		conn, err := dial(s.dsn, s.tlsConfig)
		if err != nil {
			s.logger.Error("reconnect failed",
				zap.Int("attempt", attempts),
				zap.Error(err),
			)
			delay = min(delay*_backoffMultiplier, _maxReconnectBackoff)
			continue
		}

		s.mu.Lock()
		s.conn = conn
		s.mu.Unlock()

		s.logger.Info("reconnected to rabbitmq", zap.Int("attempt", attempts))

		recreateErr := s.recreateChannels()
		if recreateErr != nil {
			s.logger.Error("failed to recreate channels after reconnect",
				zap.Error(recreateErr),
			)
			_ = conn.Close()
			delay = min(delay*_backoffMultiplier, _maxReconnectBackoff)
			continue
		}

		return
	}
}

func (s *Client) recreateChannels() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for queueName := range s.queues {
		ch, err := s.conn.Channel()
		if err != nil {
			return errors.Wrap(err, "recreate queue channel")
		}
		newQ, err := ch.QueueDeclare(
			queueName, true, false, false, false, nil,
		)
		if err != nil {
			return errors.Wrap(err, "redeclare queue")
		}
		s.queues[queueName] = newQ
		_ = ch.Close()

		prodChans := s.producers[queueName]
		for i := range prodChans {
			prodCh, prodErr := s.conn.Channel()
			if prodErr != nil {
				return errors.Wrap(prodErr, "recreate producer channel")
			}
			_ = prodChans[i].Close()
			prodChans[i] = prodCh
		}

		conss := s.consumers[queueName]
		for _, cons := range conss {
			cch, consErr := s.conn.Channel()
			if consErr != nil {
				return errors.Wrap(consErr, "recreate consumer channel")
			}
			_ = cons.ch.Close()
			cons.ch = cch

			msgs, consumeErr := cch.Consume(
				queueName, "", false, false, false, false, nil,
			)
			if consumeErr != nil {
				return errors.Wrap(consumeErr, "re-consume")
			}

			go s.consumeLoop(cons.queueName, cons, msgs, cons.processFunc)
		}
	}

	return nil
}

// DefineQueue registers a queue and launches producers and consumers.
func (s *Client) DefineQueue(
	_ context.Context,
	queueName string,
	numProducers, numConsumers int,
	processFunc func(context.Context, *Message) error,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch, err := s.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "open channel")
	}

	q, err := ch.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		return errors.Wrap(err, "declare queue")
	}
	s.queues[queueName] = q

	// Producers.
	if _, ok := s.producers[queueName]; !ok {
		s.producers[queueName] = make([]*amqp.Channel, 0, numProducers)
		for range numProducers {
			prodCh, prodErr := s.conn.Channel()
			if prodErr != nil {
				return errors.Wrap(prodErr, "producer channel")
			}
			s.producers[queueName] = append(s.producers[queueName], prodCh)
		}
		var idx uint64
		s.rrIdx.Store(queueName, &idx)
	}

	// Consumers.
	for range numConsumers {
		cch, consErr := s.conn.Channel()
		if consErr != nil {
			return errors.Wrap(consErr, "consumer channel")
		}
		cons := &consumer{
			ch:          cch,
			doneCh:      make(chan struct{}),
			queueName:   queueName,
			processFunc: processFunc,
		}
		msgs, consumeErr := cch.Consume(
			queueName, "", false, false, false, false, nil,
		)
		if consumeErr != nil {
			return errors.Wrap(consumeErr, "consume")
		}

		go s.consumeLoop(queueName, cons, msgs, processFunc)
		s.consumers[queueName] = append(s.consumers[queueName], cons)
	}

	return nil
}

func (s *Client) consumeLoop(
	queueName string,
	cons *consumer,
	msgs <-chan amqp.Delivery,
	processFunc func(context.Context, *Message) error,
) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("panic in consumer goroutine",
				zap.String("queue", queueName),
				zap.Any("recover", r),
			)
			close(cons.doneCh)
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			close(cons.doneCh)
			return
		case msg, ok := <-msgs:
			if !ok {
				close(cons.doneCh)
				return
			}

			var m Message
			unmarshalErr := json.Unmarshal(msg.Body, &m)
			if unmarshalErr != nil {
				s.logger.Error("failed to unmarshal message",
					zap.String("queue", queueName),
					zap.Error(unmarshalErr),
				)
				_ = msg.Nack(false, false)
				continue
			}

			processErr := processFunc(s.ctx, &m)
			if processErr == nil {
				_ = msg.Ack(false)
			} else {
				s.logger.Error("processFunc error",
					zap.String("queue", queueName),
					zap.Error(processErr),
				)
				s.handleRetry(&msg, queueName)
			}
		}
	}
}

func (s *Client) handleRetry(msg *amqp.Delivery, queueName string) {
	retryCount := s.getRetryCount(msg)
	if retryCount >= s.maxRetries {
		s.logger.Warn("message exceeded max retries, rejecting",
			zap.String("queue", queueName),
			zap.String("message_id", msg.MessageId),
			zap.Int("retry_count", retryCount),
		)
		_ = msg.Reject(false)

		return
	}

	_ = msg.Nack(false, true)
}

func (s *Client) getRetryCount(msg *amqp.Delivery) int {
	if msg.Headers == nil {
		return 0
	}

	deaths, ok := msg.Headers["x-death"]
	if !ok {
		return 0
	}

	deathsSlice, ok := deaths.([]any)
	if !ok || len(deathsSlice) == 0 {
		return 0
	}

	firstDeath, ok := deathsSlice[0].(amqp.Table)
	if !ok {
		return 0
	}

	count, ok := firstDeath["count"]
	if !ok {
		return 0
	}

	switch v := count.(type) {
	case int64:
		return int(v)
	case int32:
		return int(v)
	case int:
		return v
	default:
		return 0
	}
}

// SendMessage publishes a Message struct as JSON using atomic round-robin for channel selection.
func (s *Client) SendMessage(ctx context.Context, queueName string, message *Message) error {
	s.mu.Lock()
	prodChans, ok := s.producers[queueName]
	s.mu.Unlock()
	if !ok || len(prodChans) == 0 {
		return errors.New("producer not defined for queue: " + queueName)
	}

	raw, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "marshal message")
	}

	headers := amqp.Table{}
	maps.Copy(headers, message.Headers)

	val, ok := s.rrIdx.Load(queueName)
	if !ok {
		return errors.New("round-robin index not found for queue: " + queueName)
	}
	idxPtr, ok := val.(*uint64)
	if !ok {
		return errors.New("invalid round-robin index type for queue: " + queueName)
	}

	idx := atomic.AddUint64(idxPtr, 1) - 1
	//nolint:gosec // G115: idx % len(prodChans) is always within int range.
	ch := prodChans[int(idx)%len(prodChans)]

	return ch.PublishWithContext(ctx, "", queueName, false, false, amqp.Publishing{
		Body:      raw,
		MessageId: message.MessageID,
		Timestamp: message.CreatedAt,
		Headers:   headers,
	})
}

// Close gracefully closes all channels and the connection.
func (s *Client) Close() error {
	s.cancelFunc()

	s.mu.Lock()
	defer s.mu.Unlock()

	var err error

	for _, prodChans := range s.producers {
		for _, prod := range prodChans {
			e := prod.Close()
			if e != nil && err == nil {
				err = e
			}
		}
	}
	for _, conss := range s.consumers {
		for _, cons := range conss {
			_ = cons.ch.Close()
			select {
			case <-cons.doneCh:
			case <-time.After(_reconnectTimeout):
			}
		}
	}
	if s.conn != nil {
		e := s.conn.Close()
		if e != nil && err == nil {
			err = e
		}
	}

	return err
}
