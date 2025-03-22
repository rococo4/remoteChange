package k8s_deploy

import (
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"path/filepath"
)

func deploy() {
	kubeconfig := filepath.Join("/home/user/.kube", "config") // Укажите свой путь
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
}
