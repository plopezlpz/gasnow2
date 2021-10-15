package infura

type Transaction struct {
	Blockhash            string        `json:"blockHash"`
	Blocknumber          string        `json:"blockNumber"`
	From                 string        `json:"from"`
	Gas                  string        `json:"gas"`
	Gasprice             string        `json:"gasPrice"`
	Hash                 string        `json:"hash"`
	Input                string        `json:"input"`
	Nonce                string        `json:"nonce"`
	R                    string        `json:"r"`
	S                    string        `json:"s"`
	To                   string        `json:"to"`
	Transactionindex     string        `json:"transactionIndex"`
	Type                 string        `json:"type"`
	V                    string        `json:"v"`
	Value                string        `json:"value"`
	Accesslist           []interface{} `json:"accessList,omitempty"`
	Chainid              string        `json:"chainId,omitempty"`
	Maxfeepergas         string        `json:"maxFeePerGas,omitempty"`
	Maxpriorityfeepergas string        `json:"maxPriorityFeePerGas,omitempty"`
}
type Block struct {
	Basefeepergas    string        `json:"baseFeePerGas"`
	Difficulty       string        `json:"difficulty"`
	Extradata        string        `json:"extraData"`
	Gaslimit         string        `json:"gasLimit"`
	Gasused          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	Logsbloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	Mixhash          string        `json:"mixHash"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"`
	Parenthash       string        `json:"parentHash"`
	Receiptsroot     string        `json:"receiptsRoot"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	Size             string        `json:"size"`
	Stateroot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	Totaldifficulty  string        `json:"totalDifficulty"`
	Transactions     []Transaction `json:"transactions"`
	Transactionsroot string        `json:"transactionsRoot"`
	Uncles           []interface{} `json:"uncles"`
}
