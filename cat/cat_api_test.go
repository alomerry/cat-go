package cat

import (
	"context"
	"testing"
	"time"

	"github.com/alomerry/cat-go/message"
	"github.com/stretchr/testify/assert"
)

func TestNewTransactionWithCtx(t *testing.T) {
	var (
		ctx = context.TODO()
		tx  message.Transactor

		turns = 10
		tick  *time.Ticker
	)

	tick = time.NewTicker(time.Millisecond * 500)
	defer tick.Stop()

	Init("stage")
	defer Shutdown()

	t.Run("1", func(t *testing.T) {
		for i := 0; i < turns; i++ {
			tx, ctx = NewTransactionWithCtx(ctx, "Tx", "test")
			assert.NotNil(t, tx)

			tx.AddData("foo", "bar\n")
			tx.AddData("traceId", tx.TraceId()+"\n")
			tx.SetStatus(FAIL)

			tx.Complete()
			<-tick.C
		}
	})

	t.Run("2", func(t *testing.T) {
		for i := 0; i < turns; i++ {
			tx, ctx = NewTransactionWithCtx(ctx, "Tx", "test")
			assert.NotNil(t, tx)

			fn(ctx)

			tx.AddData("测试 2", "初始 tx\n")
			tx.AddData("traceId", tx.TraceId()+"\n")
			tx.SetStatus(SUCCESS)

			tx.Complete()
			<-tick.C
		}
	})
}

func fn(ctx context.Context) {
	tx, ctx := NewTransactionWithCtx(ctx, "Tx", "test2")
	defer tx.Complete()

	time.Sleep(time.Millisecond * 100)
	tx.AddData("sb", "二级 tx\n")
	tx.AddData("traceId", tx.TraceId()+"\n")
	tx.SetStatus(SUCCESS)
}
