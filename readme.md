# Kubernetes AI Analysis Tool

This tool leverages the power of OpenAI to analyze Kubernetes clusters. By fetching data from a Kubernetes cluster, such as pod statuses and events, it can diagnose potential issues and provide AI-powered suggestions for solutions.

## Features

- Fetch the status of pods in the `default` namespace.
- Retrieve events related to non-running pods.
- Use OpenAI to suggest solutions based on pod statuses and events.

## Prerequisites

- Go (recommended version X.X or higher)
- A Kubernetes cluster (with `kubectl` configured correctly)
- OpenAI API Key

## Setup

1. **Clone the Repository**:
   
   ```bash
   git clone [repository_url]
   cd [repository_name]
   ```


2. **Set Up OpenAI API Key**:

  Ensure your OpenAI API key is available as an environment variable:
   ```bash
   export OPENAI_API_KEY='your_openai_api_key'
   ```

3. **Set Up KUBECONFIG**:

  (If not using the default location ~/.kube/config):
   ```bash
   export KUBECONFIG='/path/to/your/kubeconfig'
   ```

3. **Run the tool**:
   ```bash
   go run main.go
   ```






