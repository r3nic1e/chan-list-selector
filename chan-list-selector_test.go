package chanlistselector

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type result struct {
	r1 int
	r2 interface{}
	r3 int
}

func TestChanListSelector_AddChan(t *testing.T) {
	type args struct {
		ch interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid chan",
			args: args{
				ch: make(chan string),
			},
		},
		{
			name: "invalid chan",
			args: args{
				ch: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChanListSelector{}
			err := c.AddChan(tt.args.ch)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestChanListSelector_AddChans(t *testing.T) {
	type args struct {
		chs []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty slice",
		},
		{
			name: "valid slice",
			args: args{
				chs: []interface{}{
					make(chan string),
				},
			},
		},
		{
			name: "invalid slice",
			args: args{
				chs: []interface{}{
					make(chan string),
					0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChanListSelector{}
			err := c.AddChans(tt.args.chs...)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func selectWithContext(ctx context.Context, c *ChanListSelector, ch chan result) error {
	d := make(chan bool, 1)

	go func(ctx context.Context, c *ChanListSelector, ch chan<- result, d chan<- bool) {
		r1, r2, r3 := c.Select()
		d <- true
		ch <- result{r1: r1, r2: r2, r3: r3}
	}(ctx, c, ch, d)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-d:
		return nil
	}
}

func TestChanListSelector_Empty(t *testing.T) {
	c := &ChanListSelector{}
	require.True(t, c.Empty())

	ch := make(chan string, 1)
	ch <- ""
	close(ch)

	err := c.AddChan(ch)
	require.NoError(t, err)
	require.False(t, c.Empty())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resChan := make(chan result, 1)
	err = selectWithContext(ctx, c, resChan)
	require.NoError(t, err)
	require.False(t, c.Empty())

	err = selectWithContext(ctx, c, resChan)
	require.NoError(t, err)
	require.True(t, c.Empty())
}

func TestChanListSelector_Select(t *testing.T) {
	c := &ChanListSelector{}
	require.True(t, c.Empty())

	ch := make(chan string, 1)

	err := c.AddChan(ch)
	require.NoError(t, err)
	require.False(t, c.Empty())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resChan := make(chan result, 1)
	err = selectWithContext(ctx, c, resChan)
	require.Error(t, err)
	ch <- ""
	err = selectWithContext(ctx, c, resChan)
	require.Error(t, err)
	res := <-resChan
	require.Equal(t, 0, res.r1)
	require.Equal(t, "", res.r2)
	close(ch)
	res = <-resChan
}

func Test_isChan(t *testing.T) {
	type args struct {
		c interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "channel",
			args: args{
				c: make(chan bool),
			},
			want: true,
		},
		{
			name: "int",
			args: args{
				c: 0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isChan(tt.args.c)
			require.Equal(t, got, tt.want)
		})
	}
}
