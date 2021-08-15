package request

import (
	"reflect"
	"testing"
	"time"
)

func TestRequestFrequencyLimiter(t *testing.T) {
	type args struct {
		maxRequests          int
		periodSecs           time.Duration
		sendRequests         int
		timeBetweenReqMillis time.Duration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "limit reached",
			args: args{
				maxRequests:          10,
				periodSecs:           5,
				sendRequests:         15,
				timeBetweenReqMillis: 100,
			},
			want: true,
		}, {
			name: "limit not reached",
			args: args{
				maxRequests:          10,
				periodSecs:           5,
				sendRequests:         5,
				timeBetweenReqMillis: 200,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit := FrequencyLimiter(tt.args.maxRequests, tt.args.periodSecs)

			reached := false

			for i := 0; i < tt.args.sendRequests; i++ {
				reached = limit()
				if reached {
					break
				}
				time.Sleep(tt.args.timeBetweenReqMillis * time.Millisecond)
			}

			if !reflect.DeepEqual(reached, tt.want) {
				t.Errorf("RequestFrequencyLimiter() = %v, want %v", reached, tt.want)
			}
		})
	}
}
