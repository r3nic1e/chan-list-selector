package chanlistselector

import (
	"fmt"
	"time"
)

func ExampleChanListSelector_Select() {
	ch := make(chan bool)
	t := time.NewTicker(time.Second)

	selector := &ChanListSelector{}
	if err := selector.AddChan(t); err != nil {
		panic(err)
	}
	if err := selector.AddChan(ch); err != nil {
		panic(err)
	}

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
