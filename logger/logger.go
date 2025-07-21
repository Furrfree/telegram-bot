package logger

import (
	"fmt"
	"os"
)

func Log(msgs ...any) {
	fmt.Println(msgs...)
}

func Fatal(msgs ...any) {
	Log(append([]any{" ☠️ | Fatal:"}, msgs...)...)
	os.Exit(1)
}

func Error(msgs ...any) {
	Log(append([]any{" ❌ | Error:"}, msgs...)...)
}

func Warning(msgs ...any) {
	Log(append([]any{" ⚠️  | Warning:"}, msgs...)...)
}
