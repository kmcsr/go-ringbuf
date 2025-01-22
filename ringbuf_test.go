// Ring buffer
// Copyright (C) 2025  Kevin Z <zyxkad@gmail.com>
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ringbuf_test

import (
	"testing"

	. "github.com/kmcsr/go-ringbuf"
)

func TestRingBufferPushGetPoll(t *testing.T) {
	rb := NewRingBuffer[int](3)
	expect := func(i int, v int) {
		if got := rb.Get(i); got != v {
			t.Errorf("Expect %d at i %d, got %d", v, i, got)
		}
	}

	rb.Push(1)
	expect(0, 1)
	if got, v := rb.Len(), 1; got != v {
		t.Errorf("Expect %d for length, got %d", v, got)
	}
	rb.Push(2)
	expect(0, 1)
	expect(1, 2)
	if got, v := rb.Len(), 2; got != v {
		t.Errorf("Expect %d for length, got %d", v, got)
	}
	rb.Push(3)
	expect(0, 1)
	expect(1, 2)
	expect(2, 3)
	if got, v := rb.Len(), 3; got != v {
		t.Errorf("Expect %d for length, got %d", v, got)
	}
	rb.Push(4)
	expect(0, 2)
	expect(1, 3)
	expect(2, 4)
	if got, v := rb.Len(), 3; got != v {
		t.Errorf("Expect %d for length, got %d", v, got)
	}
	rb.Push(5)
	expect(0, 3)
	expect(1, 4)
	expect(2, 5)
	rb.Push(6)
	expect(0, 4)
	expect(1, 5)
	expect(2, 6)
	if got, v := rb.Len(), 3; got != v {
		t.Errorf("Expect %d for length, got %d", v, got)
	}
	rb.Push(7)
	if got, ok := rb.Poll(); !ok || got != 5 {
		t.Errorf("Expect %d when poll, got %d", 5, got)
	}
	expect(0, 6)
	expect(1, 7)
	if got, v := rb.Len(), 2; got != v {
		t.Errorf("Expect %d for length, got %d", v, got)
	}
	rb.Push(8)
	expect(0, 6)
	expect(1, 7)
	expect(2, 8)
	if got, v := rb.Len(), 3; got != v {
		t.Errorf("Expect %d for length, got %d", v, got)
	}
	if got, ok := rb.Poll(); !ok || got != 6 {
		t.Errorf("Expect %d when poll, got %d", 6, got)
	}
	expect(0, 7)
	expect(1, 8)
	if got, v := rb.Len(), 2; got != v {
		t.Errorf("Expect %d for length, got %d", v, got)
	}
	if got, ok := rb.Poll(); !ok || got != 7 {
		t.Errorf("Expect %d when poll, got %d", 7, got)
	}
	expect(0, 8)
	if got, v := rb.Len(), 1; got != v {
		t.Errorf("Expect %d for length, got %d", v, got)
	}
}
