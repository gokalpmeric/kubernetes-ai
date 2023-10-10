package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const openaiURL = "https://api.openai.com/v1/engines/gpt-3.5-turbo-instruct/completions"

func main() {
	// Set up Kubernetes client
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Fetch pods in default namespace
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Check if any pods are not in 'Running' status
	message := "Pods in default namespace not in 'Running' status:\n"
	for _, pod := range pods.Items {
		if pod.Status.Phase != "Running" {

			// Fetch pod events
			podEvents, err := clientset.CoreV1().Events(pod.Namespace).List(context.TODO(), metav1.ListOptions{
				FieldSelector: fmt.Sprintf("involvedObject.name=%s", pod.Name),
			})
			if err != nil {
				panic(err.Error())
			}

			// Construct a detailed description of the pod and its events
			message := fmt.Sprintf("Pod %s is in %s status.\n", pod.Name, pod.Status.Phase)
			for _, event := range podEvents.Items {
				message += fmt.Sprintf("Event: %s. Reason: %s. Message: %s\n", event.Type, event.Reason, event.Message)
			}

			// Send this message to OpenAI for suggestions
			query := fmt.Sprintf("Based on the following details about my Kubernetes pod, what might be the issue and how can I resolve it? %s", message)
			response := sendToOpenAI(query)
			fmt.Println(response)
		}
	}

	// If all pods are running, update the message
	if message == "Pods in default namespace not in 'Running' status:\n" {
		message = "All pods in the default namespace are in 'Running' status."
	}

	// Send to OpenAI or print
	fmt.Println(message)
}

// ... (rest of the sendToOpenAI function if you're using it)

func sendToOpenAI(prompt string) string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	payload := map[string]interface{}{
		"prompt":     prompt,
		"max_tokens": 150, // Adjust as necessary
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err.Error())
	}

	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err.Error())
	}

	if choices, exists := response["choices"].([]interface{}); exists && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if text, ok := choice["text"].(string); ok {
				return text
			}
		}
	}
	return "No completion choices found."
}
