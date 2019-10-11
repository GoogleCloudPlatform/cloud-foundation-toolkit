package scorecard

import (
	"bufio"
	"context"
	"flag"
	"io"
	"runtime"
	"sync"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/gcv"
	"github.com/golang/glog"
)

var pipelineFlags struct {
	goRoutineCount    int
	reviewThreadCount int
	batchSize         int
	channelBufferSize int
}

func init() {
	flag.IntVar(
		&pipelineFlags.goRoutineCount,
		"pipeline.goRoutineCount",
		runtime.NumCPU(),
		"Number of goroutines for decode, and violation processing.")
	flag.IntVar(
		&pipelineFlags.reviewThreadCount,
		"pipeline.reviewThreadCount",
		2,
		"Number of goroutines for calling review with assets.")
	flag.IntVar(
		&pipelineFlags.batchSize,
		"pipeline.batchSize",
		64,
		"Number of assets to batch in review calls.")
	flag.IntVar(
		&pipelineFlags.channelBufferSize,
		"pipeline.channelBufferSize",
		runtime.NumCPU(),
		"Channel buffer size for input, decode, and violation processing.")
}

// Result is a result from the Pipeline
type Result struct {
	Violations []*validator.Violation
	Errs       []error
}

// Pipeline handles a streaming and parallel JSON unmarshal and asset review.
type Pipeline struct {
	dataInput     chan dataInput    // data input channel
	decodedInput  chan decodeResult // unmarshalled JSON from input
	reviewResults chan Result       // results from Review call
	decodeDone    sync.WaitGroup    // coordinates close(p.decodedInput)
	reviewDone    sync.WaitGroup    // coordinates close(p.reviewResults)
	validator     *gcv.Validator
}

// NewPipeline returns a new pipeline
func NewPipeline(validator *gcv.Validator) *Pipeline {
	p := &Pipeline{
		dataInput:     make(chan dataInput, pipelineFlags.channelBufferSize),
		decodedInput:  make(chan decodeResult, pipelineFlags.channelBufferSize),
		reviewResults: make(chan Result, pipelineFlags.channelBufferSize),
		validator:     validator,
	}

	glog.Infof("Starting pipeline with %d workers", pipelineFlags.goRoutineCount)
	p.decodeDone.Add(pipelineFlags.goRoutineCount)
	for i := 0; i < pipelineFlags.goRoutineCount; i++ {
		go p.decodeJson()
	}

	p.reviewDone.Add(pipelineFlags.reviewThreadCount)
	for i := 0; i < pipelineFlags.reviewThreadCount; i++ {
		go p.reviewAsset()
	}
	go runAfterDone(&p.decodeDone, func() { close(p.decodedInput) })
	go runAfterDone(&p.reviewDone, func() { close(p.reviewResults) })

	return p
}

// AddInput provides data input to the reader.  This is best performed serially
// from a goroutine which calls CloseInput when all data has been added.
func (p *Pipeline) AddInput(r io.ReadCloser) {
	var input dataInput
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		bytes := make([]byte, len(scanner.Bytes()))
		copy(bytes, scanner.Bytes())
		input.data = append(input.data, bytes)
		if pipelineFlags.batchSize <= len(input.data) {
			p.dataInput <- input
			input.data = nil
		}
	}
	if err := scanner.Err(); err != nil {
		p.dataInput <- dataInput{err: err}
	}
	_ = r.Close()
}

func runAfterDone(wg *sync.WaitGroup, f func()) {
	wg.Wait()
	f()
}

// CloseInput closes the data input channel which will
func (p *Pipeline) CloseInput() {
	glog.Info("Closing input")
	close(p.dataInput)
}

// Results returns a channel that will stream review results.
func (p *Pipeline) Results() <-chan Result {
	return p.reviewResults
}

type dataInput struct {
	data [][]byte
	err  error
}

type decodeResult struct {
	assets []*validator.Asset
	errs   []error
}

func (p *Pipeline) decodeJson() {
	for input := range p.dataInput {
		if input.err != nil {
			p.decodedInput <- decodeResult{errs: []error{input.err}}
			continue
		}

		var result decodeResult
		result.assets = make([]*validator.Asset, len(input.data))
		for idx, data := range input.data {
			pbAsset, err := getAssetFromJSON(data)
			if err != nil {
				result.errs = append(result.errs, err)
				continue
			}
			result.assets[idx] = pbAsset
		}
		p.decodedInput <- result
	}
	p.decodeDone.Done()
}

func (p *Pipeline) reviewAsset() {
	for input := range p.decodedInput {
		var result = Result{
			Errs: input.errs,
		}

		resp, err := p.validator.Review(context.Background(), &validator.ReviewRequest{
			Assets: input.assets,
		})
		if err != nil {
			result.Errs = append(result.Errs, err)
		}
		result.Violations = resp.Violations
		p.reviewResults <- result
	}
	p.reviewDone.Done()
}
