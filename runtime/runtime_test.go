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
