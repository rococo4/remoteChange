package k8s_deploy

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"remoteChange/internal/model"
	"sigs.k8s.io/yaml"
)

type KubernetesDeployer struct {
	client *kubernetes.Clientset
	repo   repo
}

func NewKubernetesDeployer(repo repo, client *kubernetes.Clientset) *KubernetesDeployer {
	return &KubernetesDeployer{client: client, repo: repo}
}

func (k *KubernetesDeployer) Deploy(entity model.ConfigEntity, toCreate bool) error {
	var raw map[string]interface{}
	if err := yaml.Unmarshal([]byte(entity.Content), &raw); err != nil {
		return fmt.Errorf("не удалось разобрать YAML: %w", err)
	}

	kind, ok := raw["kind"].(string)
	if !ok {
		return fmt.Errorf("YAML не содержит поле 'kind'")
	}

	metadata, ok := raw["metadata"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("YAML не содержит 'metadata'")
	}

	name, ok := metadata["name"].(string)
	if !ok || name == "" {
		return fmt.Errorf("YAML не содержит 'name'")
	}

	if !toCreate && name != entity.Name {
		return fmt.Errorf("название в YAML (%s) не совпадает с названием конфигурации (%s)", name, entity.Name)
	}

	namespace, ok := metadata["namespace"].(string)
	if !ok || namespace == "" {
		return fmt.Errorf("YAML не содержит 'namespace'")
	}

	isAllowed, err := k.isNamespaceAllowed(entity.TeamId, namespace)
	if err != nil {
		return fmt.Errorf("ошибка проверки namespace: %w", err)
	}

	if !isAllowed {
		return fmt.Errorf("команде %d запрещено развертывание в namespace %s", entity.TeamId, namespace)
	}

	entity.Name = name
	entity.Type = kind
	err = k.repo.SaveConfig(entity)
	if err != nil {
		return fmt.Errorf("не удалось установить имя конфигурации: %w", err)
	}

	switch kind {
	case "ConfigMap":

		return k.deployConfigMap(entity.Content)
	case "Secret":
		return k.deploySecret(entity.Content)
	default:
		return fmt.Errorf("неизвестный ресурс: %s", kind)
	}
}

// getKubeClient загружает конфигурацию Kubernetes

// deployConfigMap создаёт или обновляет ConfigMap
func (k *KubernetesDeployer) deployConfigMap(yamlContent string) error {
	var cm v1.ConfigMap
	if err := yaml.Unmarshal([]byte(yamlContent), &cm); err != nil {
		return fmt.Errorf("не удалось разобрать ConfigMap: %w", err)
	}

	cmClient := k.client.CoreV1().ConfigMaps(cm.Namespace)
	existing, err := cmClient.Get(context.TODO(), cm.Name, metav1.GetOptions{})

	if err == nil {
		// Обновляем
		cm.ResourceVersion = existing.ResourceVersion
		_, err = cmClient.Update(context.TODO(), &cm, metav1.UpdateOptions{})
	} else {
		// Создаём
		_, err = cmClient.Create(context.TODO(), &cm, metav1.CreateOptions{})
	}

	if err != nil {
		return fmt.Errorf("не удалось развернуть ConfigMap: %w", err)
	}
	return nil
}

// deploySecret создаёт или обновляет Secret
func (k *KubernetesDeployer) deploySecret(yamlContent string) error {
	var secret v1.Secret
	if err := yaml.Unmarshal([]byte(yamlContent), &secret); err != nil {
		return fmt.Errorf("не удалось разобрать Secret: %w", err)
	}

	secretClient := k.client.CoreV1().Secrets(secret.Namespace)
	existing, err := secretClient.Get(context.TODO(), secret.Name, metav1.GetOptions{})

	if err == nil {
		// Обновляем
		secret.ResourceVersion = existing.ResourceVersion
		_, err = secretClient.Update(context.TODO(), &secret, metav1.UpdateOptions{})
	} else {
		// Создаём
		_, err = secretClient.Create(context.TODO(), &secret, metav1.CreateOptions{})
	}

	if err != nil {
		return fmt.Errorf("не удалось развернуть Secret: %w", err)
	}
	return nil
}

// isNamespaceAllowed проверяет, разрешён ли namespace для teamId
func (k *KubernetesDeployer) isNamespaceAllowed(teamId int64, namespace string) (bool, error) {
	team, err := k.repo.GetTeamById(teamId)
	if err != nil {
		return false, err
	}
	return team.Namespace == namespace, nil
}
