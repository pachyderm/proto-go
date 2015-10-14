package protoprocess

import "golang.org/x/net/context"

type apiProcessor struct {
	client    Client
	apiClient APIClient
}

func newAPIProcessor(client Client, apiClient APIClient) *apiProcessor {
	return &apiProcessor{client, apiClient}
}

func (a *apiProcessor) Process(dirPath string) error {
	apiDoClient, err := a.apiClient.Do(context.Background())
	if err != nil {
		return err
	}
	return a.client.Process(dirPath, apiDoClient)
}
