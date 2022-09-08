package event

import "encoding/json"

type ExampleEvent struct {
    Desc string
    Delay int64
}

func (e *ExampleEvent) Queue() string {
   return "example"
}

func (e *ExampleEvent) Value() []byte {
   bytes, _ := json.Marshal(e)

   return bytes
}

func (e *ExampleEvent) DelaySeconds() int64 {
   return e.Delay
}

func (e *ExampleEvent) Unmarshal(value []byte) Event {
   _ = json.Unmarshal(value, &e)
   return e
}
