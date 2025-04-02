package internal

import (
	"fmt"
	"log"
	"net/http"
	"remoteChange/internal/api"
	"remoteChange/internal/domain/config"
	"remoteChange/internal/domain/team"
	"remoteChange/internal/domain/user"
	k8s_deploy "remoteChange/internal/k8s-deploy"
	"remoteChange/internal/repository"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/rs/cors"
)

func Run() {
	connStr := "user=postgres dbname=remote_changes sslmode=disable host=localhost port=5432 password=1234" // replace with your connection string
	db, err := sqlx.Open("postgres", connStr)                                                               // replace "postgres" with your driver
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
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

	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/", fs)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))

}
func getKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		configPath := "/Users/mac/.kube/config"
		config, err = clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			return nil, fmt.Errorf("не удалось загрузить конфигурацию: %w", err)
		}
	}
	return kubernetes.NewForConfig(config)
}
