package chanlistselector

import (
	"errors"
	"reflect"
)

var (
	// ErrorInvalidValue provided value is not a channel
	ErrorInvalidValue = errors.New("provided value is not a channel")
)

//ChanListSelector allows to select from multiple channels
type ChanListSelector struct {
	chans     []reflect.Value
	cases     []reflect.SelectCase
	remaining int
}

//isChan uses reflect to check whether provided value is channel
func isChan(c interface{}) bool {
	return reflect.TypeOf(c).Kind() == reflect.Chan
}

// AddChan adds single channel to select cases
// It will return error if provided value is not a channel
func (c *ChanListSelector) AddChan(ch interface{}) error {
	if !isChan(ch) {
		return ErrorInvalidValue
	}

	c.chans = append(c.chans, reflect.ValueOf(ch))
	c.cases = append(c.cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ch),
	})
	c.remaining++
	return nil
}

// AddChans adds multiple channels to select cases
// It will return error if any of provided values is not a channel
func (c *ChanListSelector) AddChans(chs ...interface{}) error {
	for _, ch := range chs {
		err := c.AddChan(ch)
		if err != nil {
			return err
		}
	}
	return nil
}

// Empty returns whether there are any non-closed channels left
func (c *ChanListSelector) Empty() bool {
	return c.remaining <= 0
}

// Select runs select for added channels and returns 3 values: index of chosen
// channel, selected value and remaining open channels.
func (c *ChanListSelector) Select() (int, interface{}, int) {
	var value interface{} = nil
	chosen := -1

	for c.remaining > 0 {
		index, v, ok := reflect.Select(c.cases)
		// If any channel is closed
		if !ok {
			// Disable case
			c.cases[index].Chan = reflect.ValueOf(nil)
			c.remaining -= 1
			continue
		} else {
			chosen = index
			value = v.Interface()
			break
		}
	}

	return chosen, value, c.remaining
}
