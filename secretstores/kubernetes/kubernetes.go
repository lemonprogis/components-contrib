// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package kubernetes

import (
	"errors"

	kubeclient "github.com/dapr/components-contrib/authentication/kubernetes"
	"github.com/dapr/components-contrib/secretstores"
	"github.com/dapr/dapr/pkg/logger"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type kubernetesSecretStore struct {
	kubeClient kubernetes.Interface
	logger     logger.Logger
}

// NewKubernetesSecretStore returns a new Kubernetes secret store
func NewKubernetesSecretStore(logger logger.Logger) secretstores.SecretStore {
	return &kubernetesSecretStore{logger: logger}
}

// Init creates a Kubernetes client
func (k *kubernetesSecretStore) Init(metadata secretstores.Metadata) error {
	client, err := kubeclient.GetKubeClient()
	if err != nil {
		return err
	}
	k.kubeClient = client

	return nil
}

// GetSecret retrieves a secret using a key and returns a map of decrypted string/string values
func (k *kubernetesSecretStore) GetSecret(req secretstores.GetSecretRequest) (secretstores.GetSecretResponse, error) {
	resp := secretstores.GetSecretResponse{
		Data: map[string]string{},
	}
	namespace, err := k.getNamespaceFromMetadata(req.Metadata)
	if err != nil {
		return resp, err
	}

	secret, err := k.kubeClient.CoreV1().Secrets(namespace).Get(req.Name, meta_v1.GetOptions{})
	if err != nil {
		return resp, err
	}

	for k, v := range secret.Data {
		resp.Data[k] = string(v)
	}

	return resp, nil
}

func (k *kubernetesSecretStore) getNamespaceFromMetadata(metadata map[string]string) (string, error) {
	if val, ok := metadata["namespace"]; ok && val != "" {
		return val, nil
	}

	return "", errors.New("namespace is missing on metadata")
}
