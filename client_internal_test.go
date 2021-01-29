package client

import (
	"reflect"
	"testing"
)

func TestExtractUnduplicateTopics(t *testing.T) {
	type args struct {
		currentTopics []string
		newTopics     []string
	}
	type want struct {
		unsubscribeTopics []string
		subscribeTopics   []string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test 01",
			args: args{
				currentTopics: []string{"/0/#"},
				newTopics:     []string{"/1/#"},
			},
			want: want{
				unsubscribeTopics: []string{"/0/#"},
				subscribeTopics:   []string{"/1/#"},
			},
		},
		{
			name: "Test 02",
			args: args{
				currentTopics: []string{"/0/#", "/1/#"},
				newTopics:     []string{"/1/#", "/2/#"},
			},
			want: want{
				unsubscribeTopics: []string{"/0/#", ""},
				subscribeTopics:   []string{"", "/2/#"},
			},
		},
		{
			name: "Test 03",
			args: args{
				currentTopics: []string{"/1/#"},
				newTopics:     []string{"/1/#"},
			},
			want: want{
				unsubscribeTopics: []string{""},
				subscribeTopics:   []string{""},
			},
		},
		{
			name: "Test 04",
			args: args{
				currentTopics: []string{"/1/2/#", "/1/3/#"},
				newTopics:     []string{"/1/#"},
			},
			want: want{
				unsubscribeTopics: []string{"/1/2/#", "/1/3/#"},
				subscribeTopics:   []string{"/1/#"},
			},
		},
		{
			name: "Test 05",
			args: args{
				currentTopics: []string{"/1/2/#", "/1/3/#"},
				newTopics:     []string{"/1/2/#", "/1/3/4/#"},
			},
			want: want{
				unsubscribeTopics: []string{"", "/1/3/#"},
				subscribeTopics:   []string{"", "/1/3/4/#"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut, st := extractUnduplicateTopics(tt.args.currentTopics, tt.args.newTopics)
			if !reflect.DeepEqual(ut, tt.want.unsubscribeTopics) {
				t.Errorf("UnsubscribeTopics Want: %v, Result: %v", tt.want.unsubscribeTopics, ut)
			}
			if !reflect.DeepEqual(st, tt.want.subscribeTopics) {
				t.Errorf("SubscribeTopics Want: %v, Result: %v", tt.want.subscribeTopics, st)
			}
		})
	}
}
