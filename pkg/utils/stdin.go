package utils

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

func StdinQuery(question string) string {

	fmt.Println(color.BlueString(question))
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}
