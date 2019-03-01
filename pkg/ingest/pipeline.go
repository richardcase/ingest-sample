package ingest

type Processor func(in <-chan interface{}, out chan interface{})

func (p Processor) Run(in <-chan interface{}) chan interface{} {
	out := make(chan interface{})
	go func() {
		p(in, out)
		close(out)
	}()
	return out
}

type Pipeline struct {
	name  string
	steps []Processor
	dest  Destination
}

func NewPipeline(name string, destination Destination, steps ...Processor) *Pipeline {
	p := &Pipeline{
		name:  name,
		dest:  destination,
		steps: steps,
	}

	return p
}

func (p *Pipeline) Run(in chan interface{}) {
	// Loop round any processing stages
	for _, step := range p.steps {
		in = step.Run(in)
	}

	// Call the destination
	p.dest(in)
}
