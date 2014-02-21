package squeue

import (
  "time"
)


type Item struct {
  Fn func()
}

type Queue struct {
  Items []Item
  Every time.Duration
}

func NewQueue(seconds time.Duration) *Queue {
  queue := new(Queue)
  queue.Items = make([]Item, 0)
  queue.Every = seconds

  return queue
}

func (q *Queue) Push(fn func()) bool {
  q.Items = append(q.Items, Item{Fn: fn})
  return true
}

func (q *Queue) next() (item Item, ok bool) {
  items := q.Items
  ok = false

  if len(items) > 0 {
    item, items = items[0], items[1:]
    q.Items = items
    ok = true
  }

  return
}

func (q *Queue) Run() {
  time.AfterFunc(q.Every, func() {
    item, ok := q.next()
    if ok {
      item.Fn()
    }
    q.Run()
  })
}
