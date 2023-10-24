/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package event

import (
	"log"
	"sync"
	"time"
)

// A fan-out event bus. Sending a event to it, it will dispatch events to all Receivers.
type FanOutBus struct {
	C           chan<- interface{} // a channel to send events
	ch          chan interface{}   // internal channel(bidirectional), C is backed by this
	recvMapLock sync.RWMutex
	recvMap     map[string]*Receiver
}

// NewFanOutBus make a new FanOutBus
//
//	bufSize: capacity of FanOutBus.C, if 0 then it's unbuffered
func NewFanOutBus(bufSize int) *FanOutBus {
	ch := make(chan interface{}, bufSize)
	return &FanOutBus{
		C:       ch,
		ch:      ch,
		recvMap: make(map[string]*Receiver, 10),
	}
}

// Make a new *Receiver and registered to event bus with fullStrategy of Drop.
// See NewRecvStrategy for more information.
func (b *FanOutBus) NewRecv(name string, bufSize int) *Receiver {
	return b.NewRecvStrategy(name, bufSize, Drop)
}

// NewRecvStrategy make a new *Receiver and registered to event bus. You can make a Receiver at any time
// except event bus is closed. Any events sent to event bus before Receiver making maybe
// or maybe not be sent to them.
//
//	name: receiver's name
//	bufSize: capacity of Receiver.C, if 0 then it's unbuffered
//	fullStrategy: if Receiver.C is full, Drop the event or Block the goroutine.
//
// Be careful if choose Block, the whole event bus will be blocked if any Receiver can't catch-up.
// i.e, if Receiver A blocks, other receivers cannot receive event until Receiver A proceeds.
func (b *FanOutBus) NewRecvStrategy(name string, bufSize int, fullStrategy FullStrategy) *Receiver {
	ch := make(chan interface{}, bufSize)

	recv := &Receiver{
		C:        ch,
		ch:       ch,
		name:     name,
		bus:      b,
		strategy: fullStrategy,
	}

	b.recvMapLock.Lock()
	b.recvMap[name] = recv
	b.recvMapLock.Unlock()

	return recv
}

// Start dispatching events
func (b *FanOutBus) GoDispatch() {
	go b.doDispatch()
}

// Close FanOutBus.C, close all Receiver.C, deregister Receivers
func (b *FanOutBus) Close() {
	close(b.ch)
	b.recvMapLock.Lock()
	for name, recv := range b.recvMap {
		close(recv.ch)
		delete(b.recvMap, name)
	}
	b.recvMapLock.Unlock()
}

func (b *FanOutBus) doDispatch() {
	for event := range b.ch {
		b.recvMapLock.RLock()
		for _, recv := range b.recvMap {
			recv.onEvent(event)
		}
		b.recvMapLock.RUnlock()
	}
}

func (b *FanOutBus) deregisterRecv(name string) {
	b.recvMapLock.Lock()
	close(b.recvMap[name].ch)
	delete(b.recvMap, name)
	b.recvMapLock.Unlock()
}

// Strategy for when Receiver.C is full
type FullStrategy int

const (
	// drop event if if Receiver.C is full
	Drop FullStrategy = iota
	// blocks until there is room for Receiver.C(buffered) or there is a receiver waiting on Receiver.C(unbuffered)
	Block FullStrategy = iota
)

// Event receiver
type Receiver struct {
	C        <-chan interface{} // a channel for receiving events
	ch       chan interface{}   // internal channel(bidirectional), C is backed by this
	name     string
	bus      *FanOutBus
	strategy FullStrategy
}

// Close Receiver.C, deregistered itself from event bus
func (r *Receiver) Close() {
	r.bus.deregisterRecv(r.name)
}

// Drain drains Receiver.C until it's empty
func (r *Receiver) Drain() (data []interface{}) {
	data = make([]interface{}, 0, 10)
drained:
	for {
		select {
		case d, ok := <-r.C:
			if !ok {
				break drained
			}
			data = append(data, d)
		default:
			break drained
		}
	}
	return data
}

// CollectTimeout collects elements of Receiver.C until timeout
func (r *Receiver) CollectTimeout(timeout time.Duration) (data []interface{}) {
	timer := time.NewTimer(timeout)
	data = make([]interface{}, 0, 10)
drained:
	for {
		select {
		case d := <-r.C:
			data = append(data, d)
		case <-timer.C:
			break drained
		}
	}
	return
}

func (r *Receiver) onEvent(v interface{}) {
	switch r.strategy {
	case Drop:
		select {
		case r.ch <- v:
		default:
			log.Println("receiver[", r.name, "]: can't catch-up, dropping event")
		}
	case Block:
		r.ch <- v
	}
}
