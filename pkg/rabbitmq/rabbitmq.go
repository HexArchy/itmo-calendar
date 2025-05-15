package rabbitmq

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Client struct {
	conn      *amqp.Connection
	queues    map[string]amqp.Queue
	consumers map[string][]*consumer
	producers map[string][]*amqp.Channel
	rrIdx     sync.Map // map[string]*uint64, round-robin index per queue.
	mu        sync.Mutex
	logger    *zap.Logger
}

type consumer struct {
	ch     *amqp.Channel
	doneCh chan struct{}
}

// New creates a new RabbitMQ client.
func New(ctx context.Context, dsn string, tls *tls.Config, logger *zap.Logger) (*Client, error) {
	var conn *amqp.Connection
	var err error
	if tls == nil {
		conn, err = amqp.Dial(dsn)
		if err != nil {
			return nil, errors.Wrap(err, "dial rabbitmq")
		}
	} else {
		conn, err = amqp.DialTLS(dsn, tls)
		if err != nil {
			return nil, errors.Wrap(err, "dial tls rabbitmq")
		}
	}

	s := &Client{
		conn:      conn,
		queues:    make(map[string]amqp.Queue),
		consumers: make(map[string][]*consumer),
		producers: make(map[string][]*amqp.Channel),
		logger:    logger,
	}

	go func() {
		<-ctx.Done()
		_ = s.Close()
	}()

	return s, nil
}

// DefineQueue registers a queue and launches producers and consumers.
func (s *Client) DefineQueue(
	ctx context.Context,
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
		for i := 0; i < numProducers; i++ {
			prodCh, err := s.conn.Channel()
			if err != nil {
				return errors.Wrap(err, "producer channel")
			}
			s.producers[queueName] = append(s.producers[queueName], prodCh)
		}
		var idx uint64
		s.rrIdx.Store(queueName, &idx)
	}

	// Consumers.
	for i := 0; i < numConsumers; i++ {
		cch, err := s.conn.Channel()
		if err != nil {
			return errors.Wrap(err, "consumer channel")
		}
		cons := &consumer{
			ch:     cch,
			doneCh: make(chan struct{}),
		}
		msgs, err := cch.Consume(
			queueName, "", false, false, false, false, nil,
		)
		if err != nil {
			return errors.Wrap(err, "consume")
		}

		go func() {
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
				case <-ctx.Done():
					close(cons.doneCh)
					return
				case msg, ok := <-msgs:
					if !ok {
						close(cons.doneCh)
						return
					}
					var m Message
					err := json.Unmarshal(msg.Body, &m)
					if err != nil {
						s.logger.Error("failed to unmarshal message",
							zap.String("queue", queueName),
							zap.Error(err),
						)
						_ = msg.Nack(false, false)
						continue
					}
					err = processFunc(ctx, &m)
					if err == nil {
						_ = msg.Ack(false)
					} else {
						s.logger.Error("processFunc error",
							zap.String("queue", queueName),
							zap.Error(err),
						)
						_ = msg.Nack(false, true)
					}
				}
			}
		}()
		s.consumers[queueName] = append(s.consumers[queueName], cons)
	}

	return nil
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
	for k, v := range message.Headers {
		headers[k] = v
	}

	val, ok := s.rrIdx.Load(queueName)
	if !ok {
		return errors.New("round-robin index not found for queue: " + queueName)
	}
	idxPtr, ok := val.(*uint64)
	if !ok {
		return errors.New("invalid round-robin index type for queue: " + queueName)
	}

	idx := atomic.AddUint64(idxPtr, 1) - 1
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
			case <-time.After(5 * time.Second):
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
