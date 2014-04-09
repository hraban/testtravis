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

package lrucache_finalizertest

import (
	"log"
	"runtime"
	"testing"
	"time"
)

// create some channels and let them go out of scope
func leakingGoroutinesHelper() {
	c := New()
	defer c.Close()
	c.Msg("foo")
	c.Msg("bar")
}

func finalizeInt(p *int) {
	// see if finalizing works at all
	log.Print("Finalizing ", *p)
}

func runFinalizer() {
	i := 3
	runtime.SetFinalizer(&i, finalizeInt)
}

func TestLeakingGoroutines(t *testing.T) {
	runFinalizer()
	time.Sleep(time.Millisecond)
	n := runtime.NumGoroutine()
	for i := 0; i < 30; i++ {
		leakingGoroutinesHelper()
	}
	// seduce the garbage collector
	starttime := time.Now()
	var x int
	for time.Since(starttime) < 5*time.Second {
		// do some crazy heap stuff
		go func(buf []byte) {
			x = copy(make([]byte, len(buf)), buf)
		}(make([]byte, 1<<16))
		runtime.Gosched()
		runtime.GC()
		time.Sleep(500 * time.Millisecond)
	}
	n2 := runtime.NumGoroutine()
	leak := n2 - n
	if leak > 0 {
		t.Error("leaking goroutines:", leak)
		panic("dumping goroutine stacks")
	}
}
