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

type Ringbuf struct {
	Input  chan<- interface{}
	Output <-chan interface{}
}

// Creates a new ring buffer of the specified size
// NB: `nil' values are ignored/filtered out
func New(size int) (result *Ringbuf) {
	input := make(chan interface{}, 0)
	output := make(chan interface{}, size)
	result = &Ringbuf{
		Input:  input,
		Output: output,
	}
	go maintain(input, output)
	return
}

// Makes sure that all writes to Input are processed before continuing
func (rb Ringbuf) Flush() {
	rb.Input <- nil
}

// maintains the input-to-output transfer
func maintain(input <-chan interface{}, output chan interface{}) {
	for {
		var n interface{} = <- input
		for n != nil {
			select {
			case output <- n:
				n = nil
			default:
				// can't write yet?
				select {
				case <-output:
					// try to remove element from output buffer
				default:
					// output buffer is empty... that reader was fast!
				}
			}
		}
	}
}
