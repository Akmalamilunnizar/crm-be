# Broadcast Follow-Up Guide

## Overview

The broadcast system now stores the list of recipient phone numbers for each broadcast, allowing you to send follow-up broadcasts to the same recipients.

## Use Case: Follow-Up After Device Recovery

**Scenario:**
1. Customer devices go down → Send "outage" broadcast to affected customers
2. Devices come back online → Send "restored" broadcast to the same customers

## Database Changes

### New Column: `recipient_phones`

The `broadcast_history` table now includes a `recipient_phones` JSON column that stores the array of phone numbers that received each broadcast.

**Migration:**
- If you already created the table, run: `database/add_recipient_phones_column.sql`
- If creating fresh, the updated `broadcast_history.sql` already includes this column

## API Endpoints

### 1. Get Recipients from Previous Broadcast

**GET `/api/broadcast/:id/recipients`**

Returns the list of phone numbers that received a specific broadcast.

**Example Request:**
```bash
GET /api/broadcast/abc123-def456-ghi789/recipients
```

**Example Response:**
```json
{
  "broadcast_id": "abc123-def456-ghi789",
  "target_group": "customers",
  "recipient_count": 5,
  "recipient_phones": [
    "081234567890",
    "081234567891",
    "081234567892",
    "081234567893",
    "081234567894"
  ],
  "sent_at": "2024-01-15T10:30:00Z",
  "message": "⚠️ INFO GANGGUAN LAYANAN ⚠️..."
}
```

### 2. Send Follow-Up Broadcast

Use the standard **POST `/api/broadcast/send`** endpoint with the phones from the previous broadcast.

**Example:**
```json
{
  "target": "customers",
  "phones": [
    "081234567890",
    "081234567891",
    "081234567892",
    "081234567893",
    "081234567894"
  ],
  "message": "✅ INFO PEMULIHAN LAYANAN ✅...",
  "template_key": "template_restored"
}
```

## Frontend Implementation Example

```typescript
// 1. Get the outage broadcast ID from history
const outageBroadcast = broadcastHistory.find(b => b.template_key === 'template_outage')

// 2. Fetch recipients from that broadcast
const recipientsResponse = await $fetch(`/api/broadcast/${outageBroadcast.id}/recipients`)

// 3. Send follow-up broadcast to same recipients
await $fetch('/api/broadcast/send', {
  method: 'POST',
  body: {
    target: 'customers',
    phones: recipientsResponse.recipient_phones,
    message: '✅ INFO PEMULIHAN LAYANAN ✅...',
    template_key: 'template_restored'
  }
})
```

## Database Schema

```sql
CREATE TABLE `broadcast_history` (
  ...
  `recipient_phones` JSON DEFAULT NULL COMMENT 'Array of phone numbers that received this broadcast',
  ...
);
```

## GORM Model

The `BroadcastHistory` model includes a `RecipientPhones` field of type `StringArray` which automatically handles JSON serialization/deserialization.

```go
type BroadcastHistory struct {
    ...
    RecipientPhones StringArray `json:"recipient_phones" gorm:"type:json"`
    ...
}
```

## Notes

1. **Team Broadcasts**: For team broadcasts, `recipient_phones` will contain all team member phone numbers from the users table at the time of sending.

2. **Customer Broadcasts**: For customer broadcasts, `recipient_phones` contains exactly the phones array sent in the request.

3. **Follow-Up Logic**: When devices come back online, you can:
   - Query the most recent "outage" broadcast for a specific template
   - Get its recipients
   - Send a "restored" broadcast to the same recipients

4. **History Response**: The `/api/broadcast/history` endpoint now includes the `id` field, so you can use it to fetch recipients.

