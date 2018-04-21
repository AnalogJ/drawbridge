package actions

import "drawbridge/pkg/config"

type DeleteAction struct {
	Config config.Interface
}

func (e *DeleteAction) All(cliAnswerDatas []map[string]interface{}, force bool) error {
	return nil
}
func (e *DeleteAction) One(cliAnswerData map[string]interface{}, force bool) error {
	return nil
}
