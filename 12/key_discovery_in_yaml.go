package main

import (
	"fmt"
	
	"github.com/go-yaml/yaml"
)

const yamlContent = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: demo
  labels:
    app: demo
spec:
  replicas: 2
  selector:
    matchLabels:
      app: demo
  template:
    metadata:
      labels:
        app: demo
    spec:
      containers:
        - name: myhello
          image: lyteabovenyte/myhello
          ports:
          - containerPort: 8888
`

func main() {
	data := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(yamlContent), &data); err != nil{
		fmt.Printf("err: %v", err)
	}

	v, ok := data["metadata"]
	if !ok {
		fmt.Println("result for key metadata not fount")
	}
	fmt.Printf("metadata: %v", v)
}

