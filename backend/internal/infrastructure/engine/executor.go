package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"smart-shift-scheduler/internal/domain"
)

type ShiftEngine struct {
	scriptPath string // Pythonファイルの場所
}

func NewShiftEngine(path string) *ShiftEngine {
	return &ShiftEngine{scriptPath: path}
}

// Generate: Pythonスクリプトを叩いてシフトを生成する
func (e *ShiftEngine) Generate(input domain.ShiftInput) (*domain.ShiftResult, error) {
	// 1. GoのstructをJSONデータに変換
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	// 2. Pythonコマンドの準備
	cmd := exec.Command("python", e.scriptPath)

	// ★追加: Pythonに「文字コードはUTF-8だぞ！」と環境変数をセットする
	// これがないとWindowsでは日本語のやり取りでエラーになります
	cmd.Env = append(cmd.Environ(), "PYTHONIOENCODING=utf-8")

	// 標準入力(Stdin)にJSONを流し込む準備
	cmd.Stdin = bytes.NewReader(inputJSON)

	// ... (後略)
	
	// 標準入力(Stdin)にJSONを流し込む準備
	cmd.Stdin = bytes.NewReader(inputJSON)

	// 標準出力(Stdout)と標準エラー(Stderr)を受け取るバッファ
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// 3. 実行！
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("python execution failed: %s, stderr: %s", err, stderr.String())
	}

	// 4. 返ってきたJSONをGoのstructに戻す
	var result domain.ShiftResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal output: %w, output: %s", err, out.String())
	}

	return &result, nil
}