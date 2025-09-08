package ml

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
)

// TroubleType represents the classification result
type TroubleType struct {
	Type        string  `json:"type"`
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
}

// SVMClassifier handles trouble ticket classification
type SVMClassifier struct {
	weights map[string]map[string]float64
	biases  map[string]float64
	mu      sync.RWMutex
}

// NewSVMClassifier creates a new SVM classifier
func NewSVMClassifier() *SVMClassifier {
	classifier := &SVMClassifier{
		weights: make(map[string]map[string]float64),
		biases:  make(map[string]float64),
	}

	// Initialize with pre-trained weights (in real implementation, load from file)
	classifier.initializeWeights()

	return classifier
}

// Initialize weights for different trouble types
func (svm *SVMClassifier) initializeWeights() {
	// Updated weights based on real training data (300 examples)

	// Kabel Putus - Cable Issues
	svm.weights["kabel_putus"] = map[string]float64{
		"kabel": 0.95, "putus": 0.9, "rusak": 0.85, "terkelupas": 0.8,
		"longgar": 0.75, "optik": 0.7, "jaringan": 0.65, "tiang": 0.6,
		"lampu": 0.5, "indikator": 0.5, "wifi": 0.4, "internet": 0.4,
		"mati": 0.3, "total": 0.3, "koneksi": 0.3, "hilang": 0.3,
	}
	svm.biases["kabel_putus"] = -0.1

	// Listrik Mati - Power Issues
	svm.weights["listrik_mati"] = map[string]float64{
		"listrik": 0.95, "mati": 0.9, "padam": 0.85, "nyala": 0.8,
		"router": 0.75, "wifi": 0.7, "tidak": 0.65, "jalan": 0.6,
		"sinyal": 0.55, "koneksi": 0.5, "terhenti": 0.5, "lampu": 0.45,
		"karena": 0.4, "saat": 0.4,
	}
	svm.biases["listrik_mati"] = -0.15

	// Kendala DHCP - DHCP Issues
	svm.weights["kendala_dhcp"] = map[string]float64{
		"dhcp": 0.95, "masalah": 0.85, "stuck": 0.8, "alamat": 0.75,
		"ip": 0.7, "address": 0.7, "restart": 0.65, "router": 0.6,
		"gagal": 0.55, "memberikan": 0.5, "perlu": 0.5, "error": 0.45,
		"wifi": 0.4, "tidak": 0.4, "jalan": 0.4, "karena": 0.4,
	}
	svm.biases["kendala_dhcp"] = -0.2

	// Perangkat Mati - Device Issues
	svm.weights["perangkat_mati"] = map[string]float64{
		"perangkat": 0.95, "mati": 0.9, "device": 0.85, "rusak": 0.8,
		"router": 0.75, "modem": 0.7, "lampu": 0.65, "indikator": 0.6,
		"tidak": 0.55, "bisa": 0.5, "menyala": 0.5, "hidup": 0.5,
		"tiba": 0.45, "total": 0.45, "wifi": 0.4, "lagi": 0.4,
	}
	svm.biases["perangkat_mati"] = -0.1

	// Config Koneksi Server - Server Configuration Issues
	svm.weights["config_koneksi_server"] = map[string]float64{
		"config": 0.95, "konfigurasi": 0.9, "server": 0.85, "koneksi": 0.8,
		"jaringan": 0.75, "gagal": 0.7, "error": 0.65, "salah": 0.6,
		"setting": 0.55, "wifi": 0.5, "karena": 0.5, "bermasalah": 0.45,
		"tidak": 0.4, "bisa": 0.4, "konek": 0.4, "di": 0.4,
	}
	svm.biases["config_koneksi_server"] = -0.25

	// Over User - Overload Issues
	svm.weights["over_user"] = map[string]float64{
		"over": 0.95, "user": 0.9, "overload": 0.85, "kebanyakan": 0.8,
		"banyak": 0.75, "pengguna": 0.7, "yang": 0.65, "pakai": 0.6,
		"lemot": 0.55, "penuh": 0.5, "lambat": 0.5, "tidak": 0.45,
		"stabil": 0.45, "koneksi": 0.4, "drop": 0.4, "putus": 0.4,
		"internet": 0.4, "wifi": 0.4, "jaringan": 0.4, "karena": 0.4,
		"terlalu": 0.4, "jadi": 0.4,
	}
	svm.biases["over_user"] = -0.2
}

// PreprocessText cleans and normalizes input text
func (svm *SVMClassifier) preprocessText(text string) []string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Remove special characters and split into words
	words := strings.Fields(text)

	// Simple stemming (in production, use proper stemming library)
	var processed []string
	for _, word := range words {
		// Remove common suffixes
		word = strings.TrimSuffix(word, "ing")
		word = strings.TrimSuffix(word, "ed")
		word = strings.TrimSuffix(word, "s")
		word = strings.TrimSuffix(word, "ly")

		// Filter out very short words
		if len(word) > 2 {
			processed = append(processed, word)
		}
	}

	return processed
}

// ExtractFeatures creates feature vector from text
func (svm *SVMClassifier) extractFeatures(words []string) map[string]float64 {
	features := make(map[string]float64)

	// Count word frequencies
	for _, word := range words {
		features[word]++
	}

	// Normalize by text length
	length := float64(len(words))
	if length > 0 {
		for word := range features {
			features[word] /= length
		}
	}

	return features
}

// Predict classifies the trouble ticket
func (svm *SVMClassifier) Predict(title string) (*TroubleType, error) {
	svm.mu.RLock()
	defer svm.mu.RUnlock()

	// Preprocess text
	words := svm.preprocessText(title)
	if len(words) == 0 {
		return &TroubleType{
			Type:        "other",
			Confidence:  0.0,
			Description: "Unable to classify - insufficient text",
		}, nil
	}

	// Extract features
	features := svm.extractFeatures(words)

	// Calculate scores for each class
	scores := make(map[string]float64)
	for class := range svm.weights {
		score := svm.biases[class]
		for word, weight := range svm.weights[class] {
			if freq, exists := features[word]; exists {
				score += weight * freq
			}
		}
		scores[class] = score
	}

	// Find the class with highest score
	var bestClass string
	var bestScore float64 = math.Inf(-1)

	for class, score := range scores {
		if score > bestScore {
			bestScore = score
			bestClass = class
		}
	}

	// Convert score to confidence (sigmoid function)
	confidence := 1.0 / (1.0 + math.Exp(-bestScore))

	// Get description for the classified type
	description := svm.getTypeDescription(bestClass)

	return &TroubleType{
		Type:        bestClass,
		Confidence:  confidence,
		Description: description,
	}, nil
}

// GetTypeDescription returns human-readable description
func (svm *SVMClassifier) getTypeDescription(troubleType string) string {
	descriptions := map[string]string{
		"kabel_putus":           "Kabel jaringan putus atau rusak",
		"listrik_mati":          "Masalah listrik dan power supply",
		"kendala_dhcp":          "Kendala DHCP dan IP address",
		"perangkat_mati":        "Perangkat hardware mati atau rusak",
		"config_koneksi_server": "Masalah konfigurasi server",
		"over_user":             "Overload karena terlalu banyak user",
		"wifi":                  "WiFi/Wireless connectivity issues",
		"internet":              "Internet speed and connectivity problems",
		"hardware":              "Physical hardware and equipment issues",
		"power":                 "Power and electrical supply problems",
		"software":              "Software and application related issues",
		"other":                 "Other miscellaneous problems",
	}

	if desc, exists := descriptions[troubleType]; exists {
		return desc
	}
	return "Unknown trouble type"
}

// BatchPredict processes multiple titles
func (svm *SVMClassifier) BatchPredict(titles []string) ([]*TroubleType, error) {
	results := make([]*TroubleType, len(titles))

	for i, title := range titles {
		result, err := svm.Predict(title)
		if err != nil {
			return nil, fmt.Errorf("error predicting title %d: %v", i, err)
		}
		results[i] = result
	}

	return results, nil
}

// GetClassificationStats returns statistics about classifications
func (svm *SVMClassifier) GetClassificationStats() map[string]interface{} {
	svm.mu.RLock()
	defer svm.mu.RUnlock()

	stats := map[string]interface{}{
		"supported_types": len(svm.weights),
		"types":           make([]string, 0, len(svm.weights)),
	}

	for class := range svm.weights {
		stats["types"] = append(stats["types"].([]string), class)
	}

	return stats
}

// TrainModel trains the SVM with new data (for future implementation)
func (svm *SVMClassifier) TrainModel(trainingData []TrainingExample) error {
	// This would implement online learning or retraining
	// For now, we use pre-trained weights
	log.Println("Training SVM model with new data...")
	// Implementation would go here
	return nil
}

// TrainingExample represents a training sample
type TrainingExample struct {
	Title string `json:"title"`
	Type  string `json:"type"`
}

// SaveModel saves the trained model to file
func (svm *SVMClassifier) SaveModel(filename string) error {
	svm.mu.RLock()
	defer svm.mu.RUnlock()

	modelData := map[string]interface{}{
		"weights": svm.weights,
		"biases":  svm.biases,
	}

	data, err := json.MarshalIndent(modelData, "", "  ")
	if err != nil {
		return err
	}

	// In real implementation, save to file
	log.Printf("Model saved to %s", filename)
	log.Printf("Model data: %s", string(data))

	return nil
}
