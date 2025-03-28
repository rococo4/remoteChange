package internal

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"remoteChange/internal/api"
	"remoteChange/internal/domain/config"
	"remoteChange/internal/domain/team"
	"remoteChange/internal/domain/user"
	k8s_deploy "remoteChange/internal/k8s-deploy"
	"remoteChange/internal/repository"
)

func Run() {
	connStr := "user=username dbname=mydb sslmode=disable" // replace with your connection string
	db, err := sqlx.Open("postgres", connStr)              // replace "postgres" with your driver
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repo := repository.NewRepo(db)

	kubeClient, err := getKubeClient()
	kubeService := k8s_deploy.NewKubernetesDeployer(repo, kubeClient)

	cfgService := config.NewService(repo, kubeService)
	teamService := team.NewService(repo)
	userService := user.NewService(repo)

	cfgHandler := api.NewCfgHandler(cfgService)
	teamHandler := api.NewTeamHandler(teamService)
	userHandler := api.NewUserHandler(userService)

	router := mux.NewRouter()
	cfgHandler.RegisterRoutes(router)
	teamHandler.RegisterTeamRoutes(router)
	userHandler.RegisterUserRoutes(router)

	log.Fatal(http.ListenAndServe(":8080", router))

}
func getKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		configPath := "/root/.kube/config"
		config, err = clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			return nil, fmt.Errorf("не удалось загрузить конфигурацию: %w", err)
		}
	}
	return kubernetes.NewForConfig(config)
}
