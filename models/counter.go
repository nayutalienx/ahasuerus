package models

type Counter struct {
	counter int
}

func NewCounter() Counter {
	return Counter{
		counter: 1,
	}
}

func (c *Counter) GetAndIncrement() int {
	temp := c.counter
	c.counter++
	return temp
}