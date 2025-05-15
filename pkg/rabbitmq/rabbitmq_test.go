package rabbitmq

import (
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func TestClient_getRetryCount(t *testing.T) {
	tests := []struct {
		name string
		msg  amqp.Delivery
		want int
	}{
		{
			name: "no headers",
			msg:  amqp.Delivery{},
			want: 0,
		},
		{
			name: "no x-death header",
			msg: amqp.Delivery{
				Headers: amqp.Table{
					"foo": "bar",
				},
			},
			want: 0,
		},
		{
			name: "x-death with count as int64",
			msg: amqp.Delivery{
				Headers: amqp.Table{
					"x-death": []any{
						amqp.Table{
							"count": int64(5),
						},
					},
				},
			},
			want: 5,
		},
		{
			name: "x-death with count as int32",
			msg: amqp.Delivery{
				Headers: amqp.Table{
					"x-death": []any{
						amqp.Table{
							"count": int32(3),
						},
					},
				},
			},
			want: 3,
		},
		{
			name: "x-death with count as int",
			msg: amqp.Delivery{
				Headers: amqp.Table{
					"x-death": []any{
						amqp.Table{
							"count": 7,
						},
					},
				},
			},
			want: 7,
		},
		{
			name: "x-death empty array",
			msg: amqp.Delivery{
				Headers: amqp.Table{
					"x-death": []any{},
				},
			},
			want: 0,
		},
	}

	logger := zap.NewNop()
	client := &Client{logger: logger}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := client.getRetryCount(&tt.msg)
			if got != tt.want {
				t.Errorf("getRetryCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
