/***************************************************************************
  Copyright (C) 2016 Christoph Reichenbach


 This program may be modified and copied freely according to the terms
 of the GNU general public license (GPL), Version 3.

 Please refer to www.gnu.org for licensing details.

 This work is provided AS IS, without warranty of any kind, expressed or
 implied, including but not limited to the warranties of merchantability,
 noninfringement, and fitness for a specific purpose. The author will not
 be held liable for any damage caused by this work or derivatives of it.

 By using this source code, you agree to the licensing terms as stated
 above.


 Please contact the maintainer for bug reports or inquiries.

 Current Maintainer:

    Christoph Reichenbach (CR) <creichen@gmail.com>

***************************************************************************/


package ringbuf

import "fmt"
import "testing"

func TestBasic(t *testing.T) {
	rb := New(2)

	for i := 0; i < 1000; i++ {
		rb.Input <- i
		if (<-rb.Output).(int) != i {
			t.Fail()
		}
	}
}

func TestSimpleOverfill(t *testing.T) {
	rb := New(2)

	rb.Input <- 0
	rb.Input <- 1
	rb.Input <- 2
	if n := (<-rb.Output).(int); n != 1 {
		t.Log(fmt.Sprintf("expected 1, got %n"), n)
		t.Fail()
	}
	if n := (<-rb.Output).(int); n != 2 {
		t.Log(fmt.Sprintf("expected 2, got %n"), n)
		t.Fail()
	}
}

func TestSimpleOverfill2(t *testing.T) {
	rb := New(2)

	rb.Input <- 0
	rb.Input <- 1
	rb.Input <- 2
	rb.Input <- 3
	rb.Flush()
	if n := (<-rb.Output).(int); n != 2 {
		t.Log(fmt.Sprintf("expected 2, got %v", n))
		t.Fail()
	}
	if n := (<-rb.Output).(int); n != 3 {
		t.Log(fmt.Sprintf("expected 3, got %v", n))
		t.Fail()
	}
}

func TestMulti(t *testing.T) {
	rb := New(2)

	rb.Input <- 0
	rb.Flush()
	if n := (<-rb.Output).(int); n != 0 {
		t.Log(fmt.Sprintf("expected 0, got %v", n))
		t.Fail()
	}
	rb.Input <- 1
	rb.Input <- 2
	rb.Input <- 3
	rb.Flush()
	if n := (<-rb.Output).(int); n != 2 {
		t.Log(fmt.Sprintf("expected 2, got %v", n))
		t.Fail()
	}
	rb.Input <- 4
	rb.Flush()
	if n := (<-rb.Output).(int); n != 3 {
		t.Log(fmt.Sprintf("expected 3, got %v", n))
		t.Fail()
	}
	if n := (<-rb.Output).(int); n != 4 {
		t.Log(fmt.Sprintf("expected 4, got %v", n))
		t.Fail()
	}
}

func TestOverfill(t *testing.T) {
	for size := 2; size < 3; size += 1 {
		for end := size; end < size*3; end++ {

			rb := New(size)
			t.Log(fmt.Sprintf("[START] ----------------------------------------\n size=%d, end=%d", size, end))
			for i := 0; i < end; i++ {
				t.Log(fmt.Sprintf("  w: %d", i))
				rb.Input <- i
			}
			rb.Flush()
			for j := 0; j < size; j++ {
				expected := j + end - size
				read := (<-rb.Output).(int)
				t.Log(fmt.Sprintf("  r: %d", read))
				if read != expected {
					t.Log(fmt.Sprintf("size=%d, end=%d, expected %d but got %d", size, end, expected, read))
					t.Fail()
				}
			}
		}
	}
}
