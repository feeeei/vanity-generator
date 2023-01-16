package runtime

import (
	"fmt"
	"log"
	"time"
	"vanity-generator/common"
	"vanity-generator/cryptos"
	"vanity-generator/cryptos/ethereum"
	"vanity-generator/cryptos/tron"
	"vanity-generator/model"
	"vanity-generator/pkg/context"
	sliding "vanity-generator/pkg/sliding_window"
)

type Executor struct {
	prefix, suffix string
	concurrency    int32
	diff           float64
	generator      cryptos.Generator

	ctx           *context.Context
	slidingWindow *sliding.SlidingWindow
	wallet        *model.Wallet
}

func NewExecutor(symbol, prefix, suffix string, concurrency int32) *Executor {
	var generator cryptos.Generator
	switch symbol {
	case "eth":
		generator = ethereum.NewEthereumGenerator()
	case "tron":
		generator = tron.NewTronGenerator()
	default:
		panic(fmt.Sprintf("unknow symbol %s", symbol))
	}

	return &Executor{
		prefix:        prefix,
		suffix:        suffix,
		concurrency:   concurrency,
		diff:          float64(generator.Difficulty(prefix, suffix)),
		ctx:           context.NewContext(),
		generator:     generator,
		slidingWindow: sliding.NewSlidingWindow(30),
	}
}

func (e *Executor) Start() *model.Wallet {
	for i := int32(0); i < e.concurrency; i++ {
		go e.exec()
	}

	go e.collection()
	go e.ShowProgress()

	e.ctx.Wait()

	return e.wallet
}

func (e *Executor) exec() {
	for !e.ctx.IsDone() {
		if wallet := e.generator.DoSingle(e.prefix, e.suffix); wallet != nil {
			e.ctx.Finish()
			e.wallet = wallet
		}
		e.ctx.Increment()
	}
}

func (e *Executor) collection() {
	ticker := time.Tick(time.Second)
	for {
		<-ticker
		power := e.ctx.Reset()
		e.slidingWindow.Add(power)
	}
}

func (e *Executor) ShowProgress() {
	start := time.Now()
	time.Sleep(time.Second + time.Millisecond*100)
	for {
		template := "duration: %v, hashing power: %s/s, expect 50%%: %-10v  70%%: %-10v  90%%: %-10v"
		duration := time.Now().Sub(start).Round(time.Second)
		power := common.HashPower(e.slidingWindow.Average())
		info := fmt.Sprintf(template, duration, power, e.P50(), e.P70(), e.P90())
		log.Println(info)

		time.Sleep(time.Second * 3)
	}
}

func (e *Executor) P50() time.Duration {
	return e.p(0.5)
}

func (e *Executor) P70() time.Duration {
	return e.p(0.7)
}

func (e *Executor) P90() time.Duration {
	return e.p(0.9)
}

func (e *Executor) p(probability float64) time.Duration {
	probability = common.Probability(probability, e.diff, e.slidingWindow.Average())

	return time.Second * time.Duration(probability)
}
