package parser

import "testing"

func TestReadBody_BulkNormal(t *testing.T) {
	// 模拟正在读取 multi bulk 的某一行：$3\r\n
	msg := []byte("$3\r\n")
	state := &readState{}

	err := readBody(msg, state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 期望：
	// bulkLen = 3
	// args = ["3"]? —— 注意：你的代码把 line ("$3") append 到 args
	if state.bulkLen != 3 {
		t.Fatalf("unexpected bulkLen: %d", state.bulkLen)
	}
	if len(state.args) != 1 {
		t.Fatalf("unexpected args count: %d", len(state.args))
	}
	if string(state.args[0]) != "$3" { // 按你的源码逻辑是 append 原始 "$3"
		t.Fatalf("unexpected arg: %q", state.args[0])
	}
}
func TestReadBody_BulkZero(t *testing.T) {
	// $0\r\n 代表空 bulk
	msg := []byte("$0\r\n")
	state := &readState{}

	err := readBody(msg, state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if state.bulkLen != 0 {
		t.Fatalf("bulkLen should be 0 for $0, got %d", state.bulkLen)
	}

	// $0 就直接 append 空 slice
	if len(state.args) != 1 {
		t.Fatalf("unexpected args count: %d", len(state.args))
	}

	if string(state.args[0]) != "" {
		t.Fatalf("expected empty bulk, got %q", state.args[0])
	}
}
func TestReadBody_InvalidNumber(t *testing.T) {
	// 非法 bulk header：例如 "$abc\r\n"
	msg := []byte("$abc\r\n")
	state := &readState{}

	err := readBody(msg, state)
	if err == nil {
		t.Fatalf("expected protocol error, got nil")
	}
}
