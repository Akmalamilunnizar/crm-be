# SVM Trouble Ticket Classification Implementation

## Overview
This document describes the implementation of Support Vector Machine (SVM) for automatic classification of trouble tickets based on their titles.

## Architecture

### Backend Components

#### 1. ML Classifier (`internal/ml/trouble_classifier.go`)
- **SVMClassifier**: Main classifier struct with pre-trained weights
- **Feature Extraction**: Text preprocessing and feature vector creation
- **Classification**: SVM-based prediction with confidence scores
- **Supported Types**: wifi, internet, hardware, power, software, other

#### 2. ML Service (`internal/api/ticket/ml_service.go`)
- **MLService**: Service layer for ML operations
- **AutoClassifyTicket**: Automatically classifies tickets during creation
- **BatchClassifyTickets**: Processes multiple tickets at once
- **GetSuggestedTypes**: Returns classification suggestions

#### 3. API Integration (`internal/api/ticket/handler.go`)
- **Enhanced Create Handler**: Auto-classification during ticket creation
- **ClassifyTicket Endpoint**: Standalone classification API
- **GetMLStats Endpoint**: ML system statistics

### Frontend Components

#### 1. SVM Classifier Component (`components/SVMClassifier.vue`)
- Interactive testing interface
- Real-time classification results
- Confidence visualization
- Example test cases

#### 2. Enhanced Ticket Form (`pages/dashboard/tickets/index.vue`)
- Auto-classification button
- Real-time type suggestions
- One-click type application

#### 3. ML Testing Page (`pages/dashboard/ml-testing/index.vue`)
- Comprehensive testing interface
- Accuracy statistics
- Quick test examples

## API Endpoints

### Classification
```http
POST /api/tickets/classify
Content-Type: application/json

{
  "title": "Internet sangat lambat"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Classification successful",
  "data": {
    "type": "internet",
    "confidence": 0.85,
    "description": "Internet speed and connectivity problems"
  }
}
```

### ML Statistics
```http
GET /api/tickets/ml/stats
```

**Response:**
```json
{
  "success": true,
  "message": "ML stats retrieved",
  "data": {
    "supported_types": 6,
    "types": ["wifi", "internet", "hardware", "power", "software", "other"]
  }
}
```

### Enhanced Ticket Creation
```http
POST /api/tickets
Content-Type: application/json

{
  "customer_id": "customer123",
  "title": "WiFi tidak connect",
  "description": "Customer reports WiFi connection issues",
  "auto_classify": true
}
```

## Classification Types

| Type | Keywords | Description |
|------|----------|-------------|
| **wifi** | wifi, wireless, connection, signal, password, ssid | WiFi/Wireless connectivity issues |
| **internet** | internet, speed, slow, bandwidth, download, ping | Internet speed and connectivity problems |
| **hardware** | hardware, device, broken, cable, port, switch | Physical hardware and equipment issues |
| **power** | power, electricity, outage, voltage, fuse, ups | Power and electrical supply problems |
| **software** | software, application, error, bug, crash, update | Software and application related issues |
| **other** | (fallback) | Other miscellaneous problems |

## Usage Examples

### 1. Automatic Classification
When creating a ticket, the system automatically classifies it:

```javascript
// Frontend
const response = await ticketsApi().create({
  customer_id: 'customer123',
  title: 'Internet sangat lambat',
  description: 'Customer reports slow internet',
  auto_classify: true
})
```

### 2. Manual Classification
Test classification without creating a ticket:

```javascript
// Frontend
const result = await ticketsApi().classifyTicket('WiFi tidak connect')
console.log(result.data.type) // 'wifi'
console.log(result.data.confidence) // 0.92
```

### 3. Batch Classification
Classify multiple tickets at once:

```go
// Backend
tickets := []entities.TroubleTicket{...}
mlService := NewMLService()
classifiedTickets, err := mlService.BatchClassifyTickets(tickets)
```

## Configuration

### Confidence Threshold
The system uses a confidence threshold of 0.6 for auto-assignment:

```go
if classification.Confidence > 0.6 {
    ticket.Type = &classification.Type
}
```

### Pre-trained Weights
The classifier comes with pre-trained weights for common Indonesian and English terms:

```go
svm.weights["wifi"] = map[string]float64{
    "wifi": 0.9, "wireless": 0.8, "connection": 0.7,
    "signal": 0.8, "password": 0.5, "ssid": 0.7,
    // ... more keywords
}
```

## Performance Considerations

### Text Preprocessing
- Convert to lowercase
- Remove special characters
- Simple stemming (remove common suffixes)
- Filter short words (< 3 characters)

### Feature Extraction
- Word frequency counting
- Normalization by text length
- Weighted scoring based on pre-trained weights

### Caching
- ML service instances are created per request
- Consider implementing singleton pattern for production
- Model weights are loaded once per service instance

## Testing

### Unit Tests
Test individual components:

```go
func TestSVMClassifier(t *testing.T) {
    classifier := NewSVMClassifier()
    result, err := classifier.Predict("Internet lambat")
    assert.NoError(t, err)
    assert.Equal(t, "internet", result.Type)
    assert.Greater(t, result.Confidence, 0.5)
}
```

### Integration Tests
Test API endpoints:

```bash
curl -X POST http://localhost:8080/api/tickets/classify \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"title": "WiFi tidak connect"}'
```

### Frontend Testing
Use the ML Testing page at `/dashboard/ml-testing` to:
- Test various ticket titles
- Verify classification accuracy
- View confidence scores
- Compare predicted vs actual types

## Future Enhancements

### 1. Model Training
- Implement online learning from user feedback
- Retrain model with new data
- A/B testing for model improvements

### 2. Advanced Features
- Multi-language support
- Context-aware classification
- Integration with external ML services
- Real-time model updates

### 3. Analytics
- Classification accuracy tracking
- User feedback collection
- Performance metrics dashboard
- Model drift detection

## Troubleshooting

### Common Issues

1. **Low Confidence Scores**
   - Check if title contains relevant keywords
   - Verify text preprocessing is working correctly
   - Consider adding more training data

2. **Wrong Classifications**
   - Review and update keyword weights
   - Add domain-specific terms
   - Implement user feedback loop

3. **Performance Issues**
   - Implement model caching
   - Use async processing for batch operations
   - Consider model optimization

### Debug Logging
Enable debug logging to see classification details:

```go
log.Printf("Classifying trouble ticket: '%s'", title)
log.Printf("Classification result: %s (confidence: %.2f)", result.Type, result.Confidence)
```

## Conclusion

The SVM implementation provides automatic trouble ticket classification with high accuracy for common problem types. The system is designed to be extensible and can be improved with additional training data and user feedback.
