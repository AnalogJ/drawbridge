package utils

import (
	"bufio"
	"os"
	"strings"
	"github.com/fatih/color"
	"fmt"
)

func StdinQuery(question string) string {

	fmt.Println(color.BlueString(question))
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}
