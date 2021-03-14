package chanlistselector

import (
	"fmt"
	"time"
)

func ExampleChanListSelector_Select() {
	ch := make(chan bool)
	t := time.NewTicker(time.Second)

	selector := &ChanListSelector{}
	selector.AddChan(t)
	selector.AddChan(ch)

	for {
		chosen, value, remaining := selector.Select()
		if selector.Empty() {
			break
		}

		if chosen == 0 {
			fmt.Printf("ticker time: %s", value.(time.Time).String())
		}
		fmt.Printf("remaining %d open channels", remaining)
	}
}
