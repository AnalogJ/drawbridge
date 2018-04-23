package actions

import "drawbridge/pkg/config"

type ProxyAction struct {
	Config config.Interface
}

func (e *ProxyAction) Start(answerDataList []map[string]interface{}, dryRun bool) error {

	// write the pac template
	pacTemplate, err := e.Config.GetPacTemplate()
	if err != nil {
		return err
	}

	_, err = pacTemplate.WriteTemplate(answerDataList, dryRun)
	if err != nil {
		return err
	}

	return nil
}
