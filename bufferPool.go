/*
 * Copyright (c) 2021 BlueStorm
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

// Package bufferPool buffer pooling technology
//
package bufferPool

import (
	"reflect"
	"sync"
	"time"
	"unsafe"
)

type ByteBuffer struct {
	CreateTime time.Time
	B          []byte
}

func (b *ByteBuffer) Len() int {
	return len(b.B)
}

func (b *ByteBuffer) Bytes() []byte {
	return b.B
}

func (b *ByteBuffer) Write(p []byte) *ByteBuffer {
	b.B = append(b.B, p...)
	return b
}

func (b *ByteBuffer) Writes(p ...byte) *ByteBuffer {
	b.B = append(b.B, p...)
	return b
}

func (b *ByteBuffer) WriteByte(c byte) *ByteBuffer {
	b.B = append(b.B, c)
	return b
}

func (b *ByteBuffer) WriteString(s string) *ByteBuffer {
	b.B = append(b.B, s...)
	return b
}

func (b *ByteBuffer) Set(p []byte) *ByteBuffer {
	b.B = append(b.B[:0], p...)
	return b
}

func (b *ByteBuffer) SetString(s string) *ByteBuffer {
	b.B = append(b.B[:0], s...)
	return b
}

func (b *ByteBuffer) String() string {
	return string(b.B)
}

func (b *ByteBuffer) StringReset() string {
	defer b.Reset()
	return string(b.B)
}

// ToString nocopy:ToString
func (b *ByteBuffer) ToString() (dst string, length int) {
	s := (*reflect.SliceHeader)(unsafe.Pointer(&b.B))
	if s.Len == 0 {
		return "", 0
	}
	d := (*reflect.StringHeader)(unsafe.Pointer(&dst))
	d.Len = s.Len
	d.Data = s.Data
	length = d.Len
	return
}

func (b *ByteBuffer) Reset() {
	b.B = b.B[:0]
}

type BufferPool struct {
	buf chan *ByteBuffer
	pm  sync.Mutex
	nm  sync.Mutex
	new int
}

func NewDefault() *BufferPool {
	return New(2000)
}

// New 创建Buffer池 max:Buffer池数量
func New(max int) *BufferPool {
	p := &BufferPool{
		buf: make(chan *ByteBuffer, max),
	}
	//初始化塞入
	for i := 0; i < max; i++ {
		p.buf <- &ByteBuffer{}
	}
	return p
}

func (p *BufferPool) New(max int) {
	buf := make(chan *ByteBuffer, max)
	//初始化塞入
	for i := 0; i < max; i++ {
		buf <- &ByteBuffer{}
	}
	//p.buf = append(p.buf, buf...)
}

func (p *BufferPool) Get() *ByteBuffer {
	select {
	case c := <-p.buf:
		return c
	default:
		p.nm.Lock()
		p.new++
		p.nm.Unlock()
		return &ByteBuffer{}
	}
}

func (p *BufferPool) Put(c *ByteBuffer) {
	p.pm.Lock()
	defer p.pm.Unlock()
	c.Reset()
	select {
	case p.buf <- c:
		break
	default:

		break
	}
}
