package connect

import "drawbridge/pkg/config"

type ConnectEngine struct {
	Config       config.Interface
}

func (e *ConnectEngine) Start(answerData map[string]interface{}) error {

}
