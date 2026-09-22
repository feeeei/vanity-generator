package args

type Args struct {
	Prefix      string
	Suffix      string
	Concurrency int32
	Mnemonic    bool  // eth、tron 是否改为从 BIP39 助记词派生
	Words       int32 // 助记词单词数，0 表示使用该链的默认值
}
