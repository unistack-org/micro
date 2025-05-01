package sql

import (
	"context"
	"database/sql/driver"
	"time"

	"go.unistack.org/micro/v4/tracer"
)

var _ driver.Tx = (*wrapperTx)(nil)

// wrapperTx defines a wrapper for driver.Tx
type wrapperTx struct {
	tx   driver.Tx
	span tracer.Span
	opts Options
	ctx  context.Context
}

// Commit implements driver.Tx Commit
func (w *wrapperTx) Commit() error {
	ts := time.Now()
	err := w.tx.Commit()
	td := time.Since(ts)
	_ = td
	if w.span != nil {
		if err != nil {
			w.span.SetStatus(tracer.SpanStatusError, err.Error())
		}
		w.span.Finish()
	}
	/*
		if w.opts.LoggerEnabled && w.opts.Logger.V(w.opts.LoggerLevel) {
			w.opts.Logger.Log(w.ctx, w.opts.LoggerLevel, w.opts.LoggerObserver(w.ctx, "Commit", getCallerName(), td, err)...)
		}
	*/
	w.ctx = nil

	return err
}

// Rollback implements driver.Tx Rollback
func (w *wrapperTx) Rollback() error {
	ts := time.Now()
	err := w.tx.Rollback()
	td := time.Since(ts)
	_ = td
	if w.span != nil {
		if err != nil {
			w.span.SetStatus(tracer.SpanStatusError, err.Error())
		}
		w.span.Finish()
	}
	/*
		if w.opts.LoggerEnabled && w.opts.Logger.V(w.opts.LoggerLevel) {
			w.opts.Logger.Log(w.ctx, w.opts.LoggerLevel, w.opts.LoggerObserver(w.ctx, "Rollback", getCallerName(), td, err)...)
		}
	*/
	w.ctx = nil

	return err
}
