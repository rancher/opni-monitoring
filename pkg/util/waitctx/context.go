package waitctx

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ttacon/chalk"
)

type waitCtxDataKeyType struct{}

var waitCtxDataKey waitCtxDataKeyType

type waitCtxData struct {
	wg sync.WaitGroup
}

func FromContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, waitCtxDataKey, &waitCtxData{
		wg: sync.WaitGroup{},
	})
}

func Background() context.Context {
	return FromContext(context.Background())
}

type RestrictiveContext = context.Context
type PermissiveContext = context.Context

type WaitContextInterface interface {
	AddOne(ctx context.Context)
	Done(ctx context.Context)
	Wait(ctx context.Context, notifyAfter ...time.Duration)

	// Go runs the given function in a background goroutine, and adds it to the
	// WaitContext. Shorthand for the following pattern:
	//  waitctx.AddOne(ctx)
	//  go func() {
	//    defer waitctx.Done(ctx)
	//    // do stuff
	//  }()
	Go(ctx context.Context, fn func())
}

type restrictive struct{}

func (restrictive) FromContext(ctx RestrictiveContext) RestrictiveContext {
	return context.WithValue(ctx, waitCtxDataKey, &waitCtxData{
		wg: sync.WaitGroup{},
	})
}

func (restrictive) AddOne(ctx RestrictiveContext) {
	data := ctx.Value(waitCtxDataKey)
	if data == nil {
		panic("context is not a WaitContext")
	}
	data.(*waitCtxData).wg.Add(1)
}

func (restrictive) Done(ctx RestrictiveContext) {
	data := ctx.Value(waitCtxDataKey)
	if data == nil {
		panic("context is not a WaitContext")
	}
	data.(*waitCtxData).wg.Done()
}

func (restrictive) Wait(ctx RestrictiveContext, notifyAfter ...time.Duration) {
	data := ctx.Value(waitCtxDataKey)
	if data == nil {
		panic("context is not a WaitContext")
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		data.(*waitCtxData).wg.Wait()
	}()
	if len(notifyAfter) > 0 {
		go func(d time.Duration) {
			for {
				select {
				case <-done:
					return
				case <-time.After(d):
					fmt.Fprint(os.Stderr, chalk.Yellow.Color("\n=== WARNING: waiting longer than expected for context to cancel ===\n"))
				}
			}
		}(notifyAfter[0])
	}
}

func (w restrictive) Go(ctx RestrictiveContext, fn func()) {
	w.AddOne(ctx)
	go func() {
		defer w.Done(ctx)
		fn()
	}()
}

type permissive struct{}

func (permissive) AddOne(ctx PermissiveContext) {
	if data := ctx.Value(waitCtxDataKey); data != nil {
		data.(*waitCtxData).wg.Add(1)
	}
}

func (permissive) Done(ctx PermissiveContext) {
	if data := ctx.Value(waitCtxDataKey); data != nil {
		data.(*waitCtxData).wg.Done()
	}
}

func (permissive) Wait(ctx PermissiveContext, notifyAfter ...time.Duration) {
	data := ctx.Value(waitCtxDataKey)
	if data == nil {
		return
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		data.(*waitCtxData).wg.Wait()
	}()
	if len(notifyAfter) > 0 {
		go func(d time.Duration) {
			for {
				select {
				case <-done:
					return
				case <-time.After(d):
					fmt.Fprint(os.Stderr, chalk.Yellow.Color("\n=== WARNING: waiting longer than expected for context to cancel ===\n"))
				}
			}
		}(notifyAfter[0])
	}
}

func (w permissive) Go(ctx PermissiveContext, fn func()) {
	w.AddOne(ctx)
	go func() {
		defer w.Done(ctx)
		fn()
	}()
}

var (
	Restrictive = restrictive{}
	Permissive  = permissive{}

	AddOne = Restrictive.AddOne
	Done   = Restrictive.Done
	Wait   = Restrictive.Wait
	Go     = Restrictive.Go
)
