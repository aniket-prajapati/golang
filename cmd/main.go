package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	//apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"

	//"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main(){
	config, err:= GetConfig()
	if err !=nil{
		return
	}


	client, err := NewClient(config)
	if err!=nil{
		return
	}

	// ------------------------------------- Deployments/Pods -------------------------------------
	//d, err := client.GetK8sClient().AppsV1().Deployments("kube-system").Get(context.Background(), "coredns", metav1.GetOptions{})
	//d, err := client.GetK8sClient().CoreV1().Pods("default").Get(context.Background(), "test-job-xvlkj", metav1.GetOptions{})
	//d, err := client.GetK8sClient().CoreV1().Pods("default").List(context.TODO(), v1.ListOptions{})
	//if err!=nil{
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(d)

	// ------------------------------------- Services -------------------------------------
	//services, err := client.GetK8sClient().CoreV1().Services("").List(context.Background(), v1.ListOptions{})
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Printf("There are %d pods in the cluster\n", len(services.Items))
	//
	//for _, s := range services.Items {
	//	for p, _ := range s.Spec.Ports {
	//		fmt.Println("Port:", s.Spec.Ports[p].Port)
	//		fmt.Println("NodePort:", s.Spec.Ports[p].NodePort)
	//	}
	//}

	// ------------------------------------- PVC -------------------------------------
	var ns, label, field string
	//var ns string
	flag.StringVar(&ns, "namespace", "default", "namespace")
	flag.StringVar(&label, "l", "", "Label selector")
	flag.StringVar(&field, "f", "", "Field selector")
	listOptions := metav1.ListOptions{
		LabelSelector: label,
		FieldSelector: field,
	}

	pvcs, err := client.GetK8sClient().CoreV1().PersistentVolumeClaims(ns).List(context.Background(), listOptions)
	if err != nil {
		log.Fatal(err)
	}
	printPVCs(pvcs)

}


func printPVCs(pvcs *v1.PersistentVolumeClaimList) {
	template := "%-32s%-8s%-8s\n"
	fmt.Printf(template, "NAME", "STATUS", "CAPACITY")
	for _, pvc := range pvcs.Items {
		quant := pvc.Spec.Resources.Requests[v1.ResourceStorage]
		fmt.Printf(
			template,
			pvc.Name,
			string(pvc.Status.Phase),
			quant.String())
	}
}

func GetConfig()(*rest.Config, error){
	var config *rest.Config
	config, err:= clientcmd.BuildConfigFromFlags("","/Users/apple/.kube/config")
	if err!=nil{
		return nil, nil
	}
	return config, nil
}


type client struct{
	k8s *kubernetes.Clientset
}

type ClientInterface interface {
	GetK8sClient() kubernetes.Interface
}

func NewClient(config *rest.Config) (ClientInterface, error){
	c:= new(client)

	k,err := kubernetes.NewForConfig(config)
	if err!=nil{
		return nil, err
	}

	c.k8s = k

	return c, err
}

func (c *client) GetK8sClient() kubernetes.Interface{
	return c.k8s
}

