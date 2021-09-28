package metrix

/*
GetPodCountsByhttp use resty package to call http url to fetch kubernetes pod count
   - TOKEN and API_URL should be part of secret created in kubernetes cluster and mounted in pod or deployment

GetPodCount method using a kubernetes rest api , in this case we need service account with cluster role permission for verb pod
to get details on cluster level.
*/
import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/common/log"
	corv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var token string

var apiUrl string

func init() {
	token = os.Getenv("TOKEN")
	apiUrl = os.Getenv("API_URL")
	//token = "f8f8a0d3bb48d67a9686c8ac0db8c41e051ecfcf9df9e65777139330f7b49f86"
	//apiUrl = "https://3b09e36a-38aa-4a4a-884e-4850f04bf71f.k8s.ondigitalocean.com"
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

func GetPodCountsByhttp() int {
	client := resty.New()
	client.SetAuthToken(token)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, err := client.R().Get(apiUrl + "/api/v1/pods")
	if err != nil {
		log.Error("http call failed", err)
	}
	podList := &corv1.PodList{}
	json.Unmarshal(resp.Body(), podList)

	totalCount := len(podList.Items)
	fmt.Println("Number of pods in cluster : ", totalCount)
	return totalCount
}
