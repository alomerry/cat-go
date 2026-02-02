package cat

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alomerry/cat-go/message"
)

func NewTransactionWithCtx(ctx context.Context, mtype, name string) (message.Transactor, context.Context) {
	if !IsEnabled() {
		return message.NullTransaction, ctx
	}

	var (
		tx message.Transactor
	)

	tx = TransactionFromCtx(ctx)
	if tx == nil {
		tx = NewTransaction(mtype, name)
	} else {
		tx = tx.NewTransaction(mtype, name)
	}

	if ctx == nil {
		ctx = context.TODO()
	}

	ctx = context.WithValue(ctx, message.CtxKeyTransaction, tx)

	return tx, ctx
}

func NewTransaction(mtype, name string) message.Transactor {
	if !IsEnabled() {
		return message.NullTransaction
	}
	return message.NewTransaction(mtype, name, manager.flush)
}

func NewCompletedTransactionWithDuration(mtype, name string, duration time.Duration) {
	if !IsEnabled() {
		return
	}

	var trans = NewTransaction(mtype, name)
	trans.SetDuration(duration)
	if duration > 0 && duration < 60*time.Second {
		trans.SetTime(time.Now().Add(-duration))
	}
	trans.SetStatus(message.CatSuccess)
	trans.Complete()
}

func NewEvent(ctx context.Context, mtype, name string) message.Messager {
	if !IsEnabled() {
		return message.NullMessage
	}

	var (
		tx message.Transactor
	)

	tx = TransactionFromCtx(ctx)
	if tx == nil {
		return message.NewEvent(mtype, name, manager.flush)
	}

	return tx.NewEvent(mtype, name)
}

func LogEvent(ctx context.Context, mtype, name string, args ...string) {
	if !IsEnabled() {
		return
	}

	var e = NewEvent(ctx, mtype, name)
	if len(args) > 0 {
		e.SetStatus(args[0])
	}
	if len(args) > 1 {
		e.SetData(args[1])
	}
	e.Complete()
}

func LogError(ctx context.Context, err error, args ...string) {
	if !IsEnabled() {
		return
	}

	var category = fmt.Sprintf("%T", err)

	if len(args) > 0 {
		category = args[0]
	}

	LogErrorWithCategory(ctx, err, category)
}

func LogErrorWithCategory(ctx context.Context, err error, category string) {
	LogErrorWithCategoryBySkipTrace(ctx, err, category, 3)
}

func LogErrorWithCategoryBySkipTrace(ctx context.Context, err error, category string, skip int, msg ...string) {
	if !IsEnabled() {
		return
	}

	if len(category) == 0 {
		category = err.Error()
	}

	var prefix string
	if len(msg) > 0 {
		prefix = strings.Join(msg, "\n") + "\n"
	}

	var event = NewEvent(ctx, "Error", category)
	var buf = newStacktrace(skip, err)
	event.SetStatus(message.CatError)
	event.SetData(prefix + buf.String())
	event.Complete()
}

func LogMetricForCount(name string, args ...int) {
	if !IsEnabled() {
		return
	}
	if len(args) == 0 {
		aggregator.metric.AddCount(name, 1)
	} else {
		aggregator.metric.AddCount(name, args[0])
	}
}

func LogMetricForDuration(name string, duration time.Duration) {
	if !IsEnabled() {
		return
	}
	aggregator.metric.AddDuration(name, duration)
}

func NewMetricHelper(name string) MetricHelper {
	if !IsEnabled() {
		return &nullMetricHelper{}
	}
	return newMetricHelper(name)
}

func TransactionFromCtx(ctx context.Context) message.Transactor {
	if ctx == nil {
		return nil
	}
	var t = ctx.Value(message.CtxKeyTransaction)
	if t == nil {
		return nil
	}
	return t.(message.Transactor)
}
