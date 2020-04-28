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
	"testing"
)

func TestFanOutBus_Close(t *testing.T) {

	t.Run("normal", func(t *testing.T) {
		bus := NewFanOutBus(1)
		bus.GoDispatch()

		foo := bus.NewRecv("foo", 0)
		bar := bus.NewRecv("bar", 0)

		bus.Close()

		if got, want := len(bus.recvMap), 0; got != want {
			t.Errorf("len(bus.recvMap) = %v, want %v", got, want)
		}

		assertChannelClosed(t, "bus.ch", bus.ch, true)
		assertChannelClosed(t, "foo.C", foo.C, true)
		assertChannelClosed(t, "bar.C", bar.C, true)

	})

	t.Run("close while receiver waiting", func(t *testing.T) {
		bus := NewFanOutBus(1)
		bus.GoDispatch()

		foo := bus.NewRecv("foo", 1)
		bar := bus.NewRecv("bar", 1)

		ready := make(chan int, 2)
		done := make(chan int, 2)

		go func() {
			ready <- 1
			<-foo.C
			assertChannelClosed(t, "foo.C", foo.C, true)
			done <- 1
		}()

		go func() {
			ready <- 1
			<-bar.C
			done <- 1
		}()

		<-ready
		<-ready

		bus.Close()

		<-done
		<-done

		if got, want := len(bus.recvMap), 0; got != want {
			t.Errorf("len(bus.recvMap) = %v, want %v", got, want)
		}

		assertChannelClosed(t, "bus.ch", bus.ch, true)
		assertChannelClosed(t, "foo.C", foo.C, true)
		assertChannelClosed(t, "bar.C", bar.C, true)

	})
}

func TestFanOutBus_NewRecv(t *testing.T) {
	bus := NewFanOutBus(1)
	foo := bus.NewRecv("foo", 1)

	if got, want := len(bus.recvMap), 1; got != want {
		t.Errorf("len(bus.recvMap) = %v, want %v", got, want)
	}

	if got, want := bus.recvMap["foo"], foo; got != want {
		t.Errorf("bus.recvMap[\"foo\"] = %v, want %v", got, want)
	}

	if got := foo.ch; got == nil {
		t.Errorf("foo.ch = nil, want not nil")
	}

	if got := foo.C; got == nil {
		t.Errorf("foo.C = nil, want not nil")
	}

	if got, want := foo.strategy, Drop; got != want {
		t.Errorf("foo.strategy = %v, want %v", got, want)
	}

	if got, want := foo.bus, bus; got != want {
		t.Errorf("foo.bus = %v, want %v", got, want)
	}
}

func TestFanOutBus_NewRecvStrategy(t *testing.T) {
	bus := NewFanOutBus(1)
	foo := bus.NewRecvStrategy("foo", 1, Block)

	if got, want := len(bus.recvMap), 1; got != want {
		t.Errorf("len(bus.recvMap) = %v, want %v", got, want)
	}

	if got, want := bus.recvMap["foo"], foo; got != want {
		t.Errorf("bus.recvMap[\"foo\"] = %v, want %v", got, want)
	}

	if got := foo.ch; got == nil {
		t.Errorf("foo.ch = nil, want not nil")
	}

	if got := foo.C; got == nil {
		t.Errorf("foo.C = nil, want not nil")
	}

	if got, want := foo.strategy, Block; got != want {
		t.Errorf("foo.strategy = %v, want %v", got, want)
	}

	if got, want := foo.bus, bus; got != want {
		t.Errorf("foo.bus = %v, want %v", got, want)
	}
}

func TestFanOutBus_GoDispatch(t *testing.T) {

	t.Run("buffered bus & receiver", func(t *testing.T) {

		bus := NewFanOutBus(1)

		foo := bus.NewRecv("foo", 1)
		bus.GoDispatch()
		// bar is registered after GoDispatch, it's ok
		bar := bus.NewRecv("bar", 1)

		bus.C <- "hello"
		if got, want := <-foo.C, "hello"; got != want {
			t.Errorf("<-foo.C = %v, want %v", got, want)
		}
		if got, want := <-bar.C, "hello"; got != want {
			t.Errorf("<-bar.C = %v, want %v", got, want)
		}

	})

	t.Run("buffered bus & unbuffered receiver", func(t *testing.T) {

		bus := NewFanOutBus(1)

		foo := bus.NewRecv("foo", 0)
		bus.GoDispatch()
		// bar is registered after GoDispatch, it's ok
		bar := bus.NewRecv("bar", 0)

		ready := make(chan int, 2)
		tested := make(chan int, 2)
		go func() {
			ready <- 1
			if got, want := <-foo.C, "hello"; got != want {
				t.Errorf("<-foo.C = %v, want %v", got, want)
			}
			tested <- 1
		}()

		go func() {
			ready <- 1
			if got, want := <-bar.C, "hello"; got != want {
				t.Errorf("<-bar.C = %v, want %v", got, want)
			}
			tested <- 1
		}()

		<-ready
		<-ready

		bus.C <- "hello"

		<-tested
		<-tested

	})

	t.Run("unbuffered bus & buffered receiver", func(t *testing.T) {
		bus := NewFanOutBus(0)

		foo := bus.NewRecv("foo", 1)
		bus.GoDispatch()
		// bar is registered after GoDispatch, it's ok
		bar := bus.NewRecv("bar", 1)

		bus.C <- "hello"

		if got, want := <-foo.C, "hello"; got != want {
			t.Errorf("<-foo.C = %v, want %v", got, want)
		}
		if got, want := <-bar.C, "hello"; got != want {
			t.Errorf("<-bar.C = %v, want %v", got, want)
		}
	})

	t.Run("unbuffered bus & unbuffered receiver", func(t *testing.T) {
		bus := NewFanOutBus(0)

		foo := bus.NewRecv("foo", 0)
		bus.GoDispatch()
		//bar is registered after GoDispatch, it's ok
		bar := bus.NewRecv("bar", 0)

		ready := make(chan int, 2)
		tested := make(chan int, 2)
		go func() {
			ready <- 1
			if got, want := <-foo.C, "hello"; got != want {
				t.Errorf("<-foo.C = %v, want %v", got, want)
			}
			tested <- 1
		}()

		go func() {
			ready <- 1
			if got, want := <-bar.C, "hello"; got != want {
				t.Errorf("<-bar.C = %v, want %v", got, want)
			}
			tested <- 1
		}()

		<-ready
		<-ready

		bus.C <- "hello"

		<-tested
		<-tested
	})

}

func TestReceiver_Close(t *testing.T) {

	t.Run("normal", func(t *testing.T) {
		bus := NewFanOutBus(1)
		bus.GoDispatch()

		foo := bus.NewRecv("foo", 0)
		bar := bus.NewRecv("bar", 0)

		foo.Close()

		if got, want := len(bus.recvMap), 1; got != want {
			t.Errorf("len(bus.recvMap) = %v, want %v", got, want)
		}

		if got, want := bus.recvMap["foo"], (*Receiver)(nil); got != want {
			t.Errorf("bus.recvMap[\"foo\"] = %v, want %v", got, want)
		}

		assertChannelClosed(t, "bus.ch", bus.ch, false)
		assertChannelClosed(t, "foo.C", foo.C, true)
		assertChannelClosed(t, "bar.C", bar.C, false)
	})

	t.Run("close while receiving", func(t *testing.T) {

		bus := NewFanOutBus(10)
		bus.GoDispatch()

		foo := bus.NewRecv("foo", 10)
		bar := bus.NewRecv("bar", 10)

		go func() {
			for i := 0; i < 10; i++ {
				bus.C <- i
			}
		}()

		done := make(chan int)
		go func() {
			foo.Close()
			close(done)
		}()

		// drain foo.C
		for i := 0; i < 10; i++ {
			<-foo.C
		}
		<-done

		bus.recvMapLock.RLock()
		if got, want := len(bus.recvMap), 1; got != want {
			t.Errorf("len(bus.recvMap) = %v, want %v", got, want)
		}
		if got, want := bus.recvMap["foo"], (*Receiver)(nil); got != want {
			t.Errorf("bus.recvMap[\"foo\"] = %v, want %v", got, want)
		}
		bus.recvMapLock.RUnlock()

		assertChannelClosed(t, "bus.ch", bus.ch, false)
		assertChannelClosed(t, "foo.C", foo.C, true)
		assertChannelClosed(t, "bar.C", bar.C, false)
	})
}

func assertChannelClosed(t *testing.T, name string, ch <-chan interface{}, expectedClosed bool) {

	if expectedClosed {
		select {
		case _, ok := <-ch:
			if ok {
				t.Errorf("<-%s is not closed", name)
			}
		default:
			t.Errorf("<-%s is not closed", name)
		}
	} else {
		select {
		case _, ok := <-ch:
			if !ok {
				t.Errorf("<-%s is closed", name)
			}
		default:
		}
	}

}
