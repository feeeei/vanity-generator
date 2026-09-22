package runtime

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"
	"vanity-generator/cmd/args"
	"vanity-generator/common"
	"vanity-generator/cryptos"
	"vanity-generator/cryptos/ethereum"
	"vanity-generator/cryptos/monero"
	"vanity-generator/cryptos/polkadot"
	"vanity-generator/cryptos/tron"
	"vanity-generator/model"
	"vanity-generator/pkg/context"
	sliding "vanity-generator/pkg/sliding_window"
)

type Executor struct {
	args      args.Args
	diff      float64
	generator cryptos.Generator

	ctx           *context.Context
	slidingWindow *sliding.SlidingWindow
	once          sync.Once
	wallet        *model.Wallet
}

func NewExecutor(symbol string, args args.Args) *Executor {
	var generator cryptos.Generator
	switch symbol {
	case "eth":
		generator = ethereum.NewEthereumGenerator(mnemonicWords(args))
	case "tron":
		generator = tron.NewTronGenerator(mnemonicWords(args))
	case "polkadot":
		generator = polkadot.NewPolkadotGenerator(wordsOrDefault(args.Words))
	case "monero":
		generator = monero.NewMoneroGenerator()
	default:
		panic(fmt.Sprintf("unknow symbol %s", symbol))
	}

	return &Executor{
		args:          args,
		diff:          float64(generator.Difficulty(args.Prefix, args.Suffix)),
		ctx:           context.NewContext(),
		generator:     generator,
		slidingWindow: sliding.NewSlidingWindow(30),
	}
}

// mnemonicWords 返回 eth、tron 使用的助记词单词数，0 表示不用助记词而直接随机私钥
func mnemonicWords(args args.Args) int {
	if !args.Mnemonic && args.Words == 0 {
		return 0
	}
	return wordsOrDefault(args.Words)
}

func wordsOrDefault(words int32) int {
	if words == 0 {
		return cryptos.DefaultMnemonicWords
	}
	return int(words)
}

func (e *Executor) Start() *model.Wallet {
	for i := int32(0); i < e.args.Concurrency; i++ {
		go e.exec()
	}

	go e.collection()
	go e.ShowProgress()

	e.ctx.Wait()

	return e.wallet
}

func (e *Executor) exec() {
	for !e.ctx.IsDone() {
		if wallet := e.generator.DoSingle(e.args.Prefix, e.args.Suffix); wallet != nil {
			// 只采用第一个结果，并且先记录结果再结束等待，避免主协程拿到 nil
			e.once.Do(func() {
				e.wallet = wallet
				e.ctx.Finish()
			})
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

// maxEtaSeconds 是 time.Duration 能表示的最长秒数，超出即视为无法估算
const maxEtaSeconds = float64(math.MaxInt64 / int64(time.Second))

func (e *Executor) ShowProgress() {
	start := time.Now()
	time.Sleep(time.Second + time.Millisecond*100)
	for {
		// 各列固定宽度，让多行输出上下对齐
		template := "duration: %-9s  hashing power: %-10s  P50: %-14s  P70: %-14s  P90: %s"
		duration := common.FormatDuration(time.Since(start))
		power := common.HashPower(e.slidingWindow.Average()) + "/s"
		info := fmt.Sprintf(template, duration, power, e.eta(0.5), e.eta(0.7), e.eta(0.9))
		log.Println(info)

		time.Sleep(time.Second * 3)
	}
}

// eta 返回以给定概率找到结果所需时长的可读形式，速度尚未统计到或结果溢出时返回 -
func (e *Executor) eta(probability float64) string {
	seconds := common.Probability(probability, e.diff, e.slidingWindow.Average())
	if math.IsNaN(seconds) || math.IsInf(seconds, 0) || seconds < 0 || seconds > maxEtaSeconds {
		return "-"
	}
	return common.FormatDuration(time.Duration(seconds * float64(time.Second)))
}
