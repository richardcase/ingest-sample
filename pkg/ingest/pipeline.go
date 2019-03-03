package ingest

// Processor presents a function that is a processing step in a pipeline
type Processor func(in <-chan interface{}, out chan interface{})

// Run will run the processing step
func (p Processor) Run(in <-chan interface{}) chan interface{} {
	out := make(chan interface{})
	go func() {
		p(in, out)
		close(out)
	}()
	return out
}

// Pipeline represents a processing pipeline that is composed of zero or more
// processing steps and a final destination
type Pipeline struct {
	name  string
	steps []Processor
	dest  Destination
}

// NewPipeline creates a new Pipeline with a name, destination and optional steps
func NewPipeline(name string, destination Destination, steps ...Processor) *Pipeline {
	p := &Pipeline{
		name:  name,
		dest:  destination,
		steps: steps,
	}

	return p
}

// Run will start the pipeline running
func (p *Pipeline) Run(in chan interface{}) {
	// Loop round any processing stages
	for _, step := range p.steps {
		in = step.Run(in)
	}

	// Call the destination
	p.dest(in)
}
