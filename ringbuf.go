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

package ringbuf

import (
	"fmt"
	"iter"
)

type RingBuffer[T any] struct {
	buf     []T
	i       int
	j       int
	hasElem bool
}

func NewRingBuffer[T any](size int) *RingBuffer[T] {
	if size < 1 {
		panic("ring buffer's size must be greater than 0")
	}
	return &RingBuffer[T]{
		buf:     make([]T, size),
		i:       0,
		j:       0,
		hasElem: false,
	}
}

// Push puts an element into the ring buffer
// It will overwrite the earliest element if there is no space avaliable
func (r *RingBuffer[T]) Push(v T) {
	r.buf[r.j] = v
	if r.hasElem {
		if r.j == r.i {
			r.i++
			if r.i == len(r.buf) {
				r.i = 0
			}
		}
	} else {
		r.hasElem = true
	}
	r.j++
	if r.j == len(r.buf) {
		r.j = 0
	}
}

// Poll removes the earliest pushed element from the ring buffer
func (r *RingBuffer[T]) Poll() (v T, ok bool) {
	if !r.hasElem {
		return v, false
	}
	v, r.buf[r.i] = r.buf[r.i], v
	r.i++
	if r.i == len(r.buf) {
		r.i = 0
	}
	if r.i == r.j {
		r.hasElem = false
	}
	return v, true
}

// Get returns the i-th element in the buffer
// It will panic if index is out of bounds
func (r *RingBuffer[T]) Get(index int) T {
	if !r.hasElem {
		panic(fmt.Errorf("Index %d out of bounds: buffer is empty", index))
	}
	if index < 0 {
		panic(fmt.Errorf("Index %d out of bounds", index))
	}
	i := index + r.i
	if r.i >= r.j {
		if i >= len(r.buf) {
			i -= len(r.buf)
			if i >= r.j {
				panic(fmt.Errorf("Index %d out of bounds", index))
			}
		}
	} else if i >= r.j {
		panic(fmt.Errorf("Index %d out of bounds", index))
	}
	return r.buf[i]
}

// Len returns the used space of the buffer
func (r *RingBuffer[T]) Len() int {
	if !r.hasElem {
		return 0
	}
	n := r.j - r.i
	if n <= 0 {
		n += len(r.buf)
	}
	return n
}

// Cap returns the total space of the buffer
func (r *RingBuffer[T]) Cap() int {
	return len(r.buf)
}

// Clear set ring buffer's length to zero
// It does not dereference old elements
func (r *RingBuffer[T]) Clear() {
	r.i = 0
	r.j = 0
	r.hasElem = false
}

// Reset set ring buffer's length to zero and dereference all elements
func (r *RingBuffer[T]) Reset() {
	r.i = 0
	r.j = 0
	r.hasElem = false
	var empty T
	for i := range len(r.buf) {
		r.buf[i] = empty
	}
}

// ForEach iterate the buffer from first to last
// if the iterator returns false, the iterate will break
func (r *RingBuffer[T]) ForEach(iter func(v T) bool) {
	if !r.hasElem {
		return
	}
	if r.j > r.i {
		for i := r.i; i < r.j; i++ {
			if !iter(r.buf[i]) {
				return
			}
		}
		return
	}
	for i := r.i; i < len(r.buf); i++ {
		if !iter(r.buf[i]) {
			return
		}
	}
	for i := 0; i < r.j; i++ {
		if !iter(r.buf[i]) {
			return
		}
	}
}

// ForEachReversed iterate the buffer from last to first
// if the iterator returns false, the iterate will break
func (r *RingBuffer[T]) ForEachReversed(iter func(v T) bool) {
	if !r.hasElem {
		return
	}
	if r.j > r.i {
		for i := r.j - 1; i >= r.i; i-- {
			if !iter(r.buf[i]) {
				return
			}
		}
		return
	}
	for i := r.j - 1; i >= 0; i-- {
		if !iter(r.buf[i]) {
			return
		}
	}
	for i := len(r.buf) - 1; i >= r.i; i-- {
		if !iter(r.buf[i]) {
			return
		}
	}
}

// Iter returns an iterator of the buffer that iterate from first to last
func (r *RingBuffer[T]) Iter() iter.Seq[T] {
	return r.ForEach
}

// IterReversed returns an iterator of the buffer that iterate from last to first
func (r *RingBuffer[T]) IterReversed() iter.Seq[T] {
	return r.ForEachReversed
}
