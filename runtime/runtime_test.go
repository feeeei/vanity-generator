package runtime

import (
	"testing"

	"vanity-generator/cmd/args"
)

// 前缀极易匹配时多个 worker 几乎同时出结果，Start 必须稳定返回非空钱包
func TestStartReturnsWallet(t *testing.T) {
	for i := 0; i < 20; i++ {
		wallet := NewExecutor("monero", args.Args{Prefix: "4", Concurrency: 8}).Start()
		if wallet == nil || wallet.Public == "" || wallet.Mnemonic == "" {
			t.Fatalf("round %d: unexpected wallet %+v", i, wallet)
		}
	}
}

func TestEta(t *testing.T) {
	// 0x 后 4 位十六进制数字，难度 16^4 = 65536
	e := NewExecutor("eth", args.Args{Prefix: "0x1234", Concurrency: 1})
	if got := e.eta(0.5); got != "-" {
		t.Errorf("eta without any speed sample = %s, want -", got)
	}
	// 每秒 1000 次时，50% 概率需要 floor(ln0.5 / ln(1-1/65536)) = 45425 次，约 45 秒
	e.slidingWindow.Add(1000)
	if got := e.eta(0.5); got != "45s" {
		t.Errorf("eta(0.5) = %s, want 45s", got)
	}
	// 每秒 1 次时约 12.6 小时
	e = NewExecutor("eth", args.Args{Prefix: "0x1234", Concurrency: 1})
	e.slidingWindow.Add(1)
	if got := e.eta(0.5); got != "12h37m5s" {
		t.Errorf("eta(0.5) at 1/s = %s, want 12h37m5s", got)
	}
}
