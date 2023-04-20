package promewrite

import (
	"testing"
"time"
)

func TestRemoteWrite(t *testing.T) {
	c, err := NewClient("http://localhost:9090/api/v1/write", 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	metrics := []MetricPoint{
		{Metric: "opcai1",
			TagsMap: map[string]string{"env": "testing", "op": "opcai"},
			Time:    time.Now().Add(-1 * time.Minute).Unix(),
			Value:   1},
		{Metric: "opcai2",
			TagsMap: map[string]string{"env": "testing", "op": "opcai"},
			Time:    time.Now().Add(-2 * time.Minute).Unix(),
			Value:   2},
		{Metric: "opcai3",
			TagsMap: map[string]string{"env": "testing", "op": "opcai"},
			Time:    time.Now().Unix(),
			Value:   3},
		{Metric: "opcai4",
			TagsMap: map[string]string{"env": "testing", "op": "opcai"},
			Time:    time.Now().Unix(),
			Value:   4},
	}
	err = c.RemoteWrite(metrics)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("end...")
}
