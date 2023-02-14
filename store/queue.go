package store

import (
	"container/list"
	"context"

	"github.com/pkg/errors"

	"github.com/floorstyle/fifo-queue/model"
	"github.com/floorstyle/fifo-queue/util"
)

type FifoQueueReal struct {
	MsgStore *list.List
}

func NewFifoQueue() model.FifoQueue {
	return &FifoQueueReal{list.New()}
}

func (queue *FifoQueueReal) PushFrontMsg(ctx context.Context, input interface{}) error {
	queue.MsgStore.PushFront(input)
	return nil
}

func (queue *FifoQueueReal) PushBackMsg(ctx context.Context, input interface{}) error {
	queue.MsgStore.PushBack(input)
	return nil
}

func (queue *FifoQueueReal) GetBackMsg(ctx context.Context, out interface{}) error {
	var msg list.Element
	err := queue.getTailMsg(ctx, &msg)
	if err != nil {
		return err
	}
	if msg.Value == nil {
		return nil
	}
	return util.UnMarshalStruct(msg.Value, out)
}

func (queue *FifoQueueReal) PopBackMsg(ctx context.Context, out interface{}) error {
	var msg list.Element
	err := queue.getTailMsg(ctx, &msg)
	if err != nil {
		return err
	}
	queue.MsgStore.Remove(&msg)
	if msg.Value == nil {
		return nil
	}
	return util.UnMarshalStruct(msg.Value, out)
}

func (queue *FifoQueueReal) GetMsgs(ctx context.Context, out interface{}) error {
	var msgs []interface{}

	msg := queue.MsgStore.Back()
	for msg != nil {
		var res interface{}
		if msg.Value != nil {
			util.UnMarshalStruct(msg.Value, &res)
		}
		msgs = append(msgs, res)
		msg = msg.Prev()
	}
	return util.UnMarshalStruct(msgs, out)
}

func (queue *FifoQueueReal) getTailMsg(ctx context.Context, out *list.Element) error {
	if queue.MsgStore.Len() == 0 {
		return errors.New(util.QUEUE_ERROR_IS_EMPTY)
	}
	msg := queue.MsgStore.Back()
	if msg == nil {
		return errors.New("stored msg is empty")
	}
	if out != nil {
		*out = *msg
	}
	return nil
}
