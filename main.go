package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
)

// Joke structure
type Joke struct {
	Message string `json:"message"`
}

// Madlib structure
type Madlib struct {
	Message string `json:"message"`
}

// EC2Response defines the structure for EC2 response
type EC2Response struct {
	Instances []string `json:"instances"`
}

// VPCResponse defines the structure for VPC response
type VPCResponse struct {
	VPCs []string `json:"vpcs"`
}

// EKSResponse defines the structure for EKS response
type EKSResponse struct {
	Clusters []string `json:"clusters"`
}

// QuotaResponse defines the structure for Quotas
type QuotaResponse struct {
	Service string `json:"service"`
	Quota   string `json:"quota"`
}

// Slice of jokes
var jokes = []string{
	"What do space cows say? Mooooooooon.",
	"What do you call a cow during an earth quake? A milk shake!",
	"What happens to an illegally parked frog? It gets toad away!",
}

// Slice of madlib components
var names = []string{"PY", "Matt", "Kia", "Shama", "Aj"}
var occupations = []string{"devOp Eng", "doctor", "pharmacist", "politician", "actor"}
var devices = []string{"laptop", "tablet", "smartphone", "desktop", "smartwatch"}
var bodyParts = []string{"wrist", "neck", "ankle", "thigh", "shoulder"}
var moods = []string{"happy", "sad", "anxious", "excited", "angry"}
var actions = []string{"playing calming music", "displaying motivational images", "showing a quick joke", "vibrating", "sending a life quote"}

// Function to get a random joke
func getJoke() string {
	rand.Seed(time.Now().UnixNano())
	return jokes[rand.Intn(len(jokes))]
}

// Function to get a random madlib
func getMadlib() string {
	rand.Seed(time.Now().UnixNano())
	name := names[rand.Intn(len(names))]
	age := rand.Intn(40) + 18
	occupation := occupations[rand.Intn(len(occupations))]
	device := devices[rand.Intn(len(devices))]
	bodyPart := bodyParts[rand.Intn(len(bodyParts))]
	mood := moods[rand.Intn(len(moods))]
	action := actions[rand.Intn(len(actions))]

	return name + " is a " + strconv.Itoa(age) + "-year-old " + occupation + " who has been struggling with a lot of job-related stress. He/she decided to try a new application to relieve stress, which runs on a/an " + device + " to help improve his/her mood. The application senses his/her mood through a device he/she wears on his/her " + bodyPart + ". When the device senses that he/she is " + mood + ", it responds by " + action + "."
}

// JokeHandler returns a random joke
func JokeHandler(w http.ResponseWriter, r *http.Request) {
	joke := getJoke()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Joke{Message: joke})
}

// MadlibHandler returns a random madlib
func MadlibHandler(w http.ResponseWriter, r *http.Request) {
	madlib := getMadlib()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Madlib{Message: madlib})
}

// HealthCheckHandler returns 200 without any logic
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Fetch EC2 Instances with error handling
func getEC2Instances() ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %v", err)
	}

	svc := ec2.NewFromConfig(cfg)

	output, err := svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to describe EC2 instances: %v", err)
	}

	var instanceIDs []string
	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, aws.ToString(instance.InstanceId))
		}
	}
	return instanceIDs, nil
}

// Fetch VPCs with error handling
func getVPCs() ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %v", err)
	}

	svc := ec2.NewFromConfig(cfg)

	output, err := svc.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to describe VPCs: %v", err)
	}

	var vpcIDs []string
	for _, vpc := range output.Vpcs {
		vpcIDs = append(vpcIDs, aws.ToString(vpc.VpcId))
	}
	return vpcIDs, nil
}

// Fetch EKS Clusters with error handling
func getEKSClusters() ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %v", err)
	}

	svc := eks.NewFromConfig(cfg)

	output, err := svc.ListClusters(context.TODO(), &eks.ListClustersInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list EKS clusters: %v", err)
	}

	return output.Clusters, nil
}

// Fetch Quotas with error handling
func getServiceQuotas(serviceCode string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return "", fmt.Errorf("unable to load AWS config: %v", err)
	}

	svc := servicequotas.NewFromConfig(cfg)

	output, err := svc.GetServiceQuota(context.TODO(), &servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String("L-1216C47A"), // Example quota code
		ServiceCode: aws.String(serviceCode),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get service quota: %v", err)
	}

	return fmt.Sprintf("%f", *output.Quota.Value), nil
}

// AWS Handlers with improved error responses
func EC2Handler(w http.ResponseWriter, r *http.Request) {
	instances, err := getEC2Instances()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := EC2Response{Instances: instances}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func VPCHandler(w http.ResponseWriter, r *http.Request) {
	vpcs, err := getVPCs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := VPCResponse{VPCs: vpcs}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func EKSHandler(w http.ResponseWriter, r *http.Request) {
	clusters, err := getEKSClusters()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := EKSResponse{Clusters: clusters}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func QuotaHandler(w http.ResponseWriter, r *http.Request) {
	quota, err := getServiceQuotas("ec2") // Example for EC2 quota
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := QuotaResponse{Service: "EC2", Quota: quota}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Jokes and Madlibs routes
	http.HandleFunc("/joke", JokeHandler)
	http.HandleFunc("/madlib", MadlibHandler)

	// AWS service routes
	http.HandleFunc("/ec2s", EC2Handler)
	http.HandleFunc("/vpcs", VPCHandler)
	http.HandleFunc("/eks", EKSHandler)
	http.HandleFunc("/quotas", QuotaHandler)

	// Health check route
	http.HandleFunc("/health", HealthCheckHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
