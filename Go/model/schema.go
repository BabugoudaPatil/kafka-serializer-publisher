package model

var _ Publishable = &JsonRequest{}

// JsonRequest - implements Publishable
type JsonRequest struct {
	Id      string            `json:"id" validate:"-"`
	Channel string            `json:"topic" validate:"required"`
	Body    interface{}       `json:"payload" validate:"required" swaggertype:"object"`
	Heads   map[string]string `json:"headers" validate:"required"`
}

func (j *JsonRequest) ID() string {
	return j.Id
}

func (j *JsonRequest) Headers() map[string]string {
	return j.Heads
}

func (j *JsonRequest) Payload() interface{} {
	return j.Body
}

func (j *JsonRequest) Topic() string {
	return j.Channel
}

func (j *JsonRequest) Subject() *string {
	return nil
}

var _ Publishable = &AvroRequest{}

// AvroRequest - implements Publishable
type AvroRequest struct {
	AvroSource string `json:"avroSource" validate:"required"`
	JsonRequest
}

func (j *AvroRequest) Subject() *string {
	return &j.AvroSource
}
