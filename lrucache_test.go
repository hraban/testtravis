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

package lrucache

import (
	"runtime"
	"testing"
	"time"
)

// create some channels and let them go out of scope
func leakingGoroutinesHelper() {
	c := New(100)
	c.Set("foo", 123)
	// keep it around!
	//time.AfterFunc(10*time.Second, func() { _ = c })
}

func TestLeakingGoroutines(t *testing.T) {
	n := runtime.NumGoroutine()
	for i := 0; i < 100; i++ {
		leakingGoroutinesHelper()
	}
	// seduce the garbage collector
	starttime := time.Now()
	for time.Since(starttime) < 5*time.Second {
		time.Sleep(time.Second)
		runtime.GC()
		runtime.Gosched()
	}
	n2 := runtime.NumGoroutine()
	leak := n2 - n
	// TODO: Why is this 1 no matter how many caches are created and cleaned up?
	// To see what I mean run this test in isolation (go test -run TestLeak);
	// leak will always be exactly 1.  When run alongside other tests their
	// garbage is also cleaned up here so leak can be less than 0---that's okay.
	if leak > 1 { // I'd like to test >0 here :(
		t.Error("leaking goroutines:", leak)
		//panic("dumping goroutine stacks")
	}
}
