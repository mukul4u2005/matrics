package matrix

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var token string

var apiUrl string

func init() {
	token = os.Getenv("TOKEN")
	apiUrl = os.Getenv("API_URL")
}

func GetPodCount(name string) int {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := client.CoreV1().Pods(name).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%d pods running in cluster\n", len(pods.Items))
	return len(pods.Items)
}
