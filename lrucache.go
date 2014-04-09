// Copyright © 2012, 2013 Lrucache contributors, see AUTHORS file
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

// Reduced code to test finalizers for lrucache project
//
// See https://github.com/hraban/lrucache
package lrucache_finalizertest

import (
	"log"
)

type Cache struct {
	opCount int
	// Cache operations are pushed down this channel to the main cache loop
	opChan chan operation
}

// Requests that are passed to the cache managing goroutine
type operation struct {
	// Cache is explicitly passed with every request so the main loop does not
	// need to keep a reference to the cache around. This allows garbage
	// collection to kick in (and eventually end the main loop through a
	// finalizer) when the last external reference to the cache goes out of
	// scope.
	c   *Cache
	msg string
}

// Consume an operation from the channel and process it. Returns false if the
// channel was closed and the main loop should stop.
//
// Implemented as a separate function to ensure all local variables go out of
// scope when this main loop iteration is complete.
//
// Imagine this function accepted an operation directly and the mainLoop
// function were implemented as follows:
//
//     for op := range opchan {
//         mainLoopBody(op)
//     }
//
// This blocks on the read from opchan, but it is not immediately clear if the
// operation from the last iteration (haha) is cleared / garbage collected while
// this read is blocking. Because the operation struct contains a reference to
// the Cache, if that doesn't happen the entire cache will not be garbage
// collected.
func mainLoopBody(opchan <-chan operation) bool {
	op, ok := <-opchan
	if !ok {
		return false
	}
	c := op.c
	c.opCount++
	log.Printf("Message #%d: %s", c.opCount, op.msg)
	return true
}

// does not keep any reference to the cache so it can be garbage collected
func mainLoop(opchan <-chan operation) {
	for mainLoopBody(opchan) {
	}
}

func (c *Cache) Init() {
	c.opChan = make(chan operation)
	go mainLoop(c.opChan)
	return
}

func (c *Cache) Msg(msg string) {
	c.opChan <- operation{c, msg}
}

func finalizeCache(c *Cache) {
	close(c.opChan)
}

func (c *Cache) Close() error {
	finalizeCache(c)
	return nil
}

// Create and initialize a new cache, ready for use.
func New() *Cache {
	var mem Cache
	c := &mem
	c.Init()
	//runtime.SetFinalizer(c, finalizeCache)
	return c
}
