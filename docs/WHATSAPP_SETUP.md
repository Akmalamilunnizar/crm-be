# Meta WhatsApp Business API Setup Guide

This guide will help you set up Meta WhatsApp Business API to send real WhatsApp messages from your application.

## Prerequisites

1. **Meta Business Account** - You need a Meta Business account
2. **WhatsApp Business Account** - Your business must have a WhatsApp Business account
3. **Phone Number** - A phone number that can receive SMS for verification

## Step 1: Create Meta App

1. Go to [Meta for Developers](https://developers.facebook.com/)
2. Click "My Apps" → "Create App"
3. Select "Business" as the app type
4. Fill in your app details:
   - App Name: Your business name
   - App Contact Email: Your email
   - Business Account: Select your business account

## Step 2: Add WhatsApp Product

1. In your app dashboard, click "Add Product"
2. Find "WhatsApp" and click "Set up"
3. You'll see the WhatsApp API setup page

## Step 3: Get Your Credentials

### Phone Number ID
1. In WhatsApp API setup, you'll see "Phone number ID"
2. Copy this ID (it looks like: `123456789012345`)

### Access Token
1. In the same page, find "Temporary access token"
2. Copy this token (it starts with `EAAB...`)
3. **Note**: This is temporary. For production, you need a permanent token.

### API URL
- Use: `https://graph.facebook.com/v18.0` (or latest version)

## Step 4: Configure Environment Variables

Create a `.env` file in your project root with:

```bash
# Meta WhatsApp Business API Configuration
WHATSAPP_API_URL=https://graph.facebook.com/v18.0
WHATSAPP_API_TOKEN=your_meta_access_token_here
WHATSAPP_PHONE_ID=your_phone_number_id_here
```

## Step 5: Test the Integration

1. Restart your Go application
2. Send a test message using the API:

```bash
curl -X POST http://localhost:3001/send-message \
  -H "Content-Type: application/json" \
  -d '{"to":"628970833338","message":"Hello from WhatsApp Business API!"}'
```

## Phone Number Format

- Use format: `628123456789` (without +)
- Include country code (62 for Indonesia)
- No spaces or special characters

## Production Setup

### Permanent Access Token
1. Go to Meta Business Manager
2. Navigate to WhatsApp → API Setup
3. Generate a permanent access token
4. Replace the temporary token in your environment variables

### Webhook Setup (Optional)
If you want to receive incoming messages:

1. Set up a webhook URL in your Meta app
2. Add webhook verification token to your environment:
   ```bash
   WHATSAPP_WEBHOOK_VERIFY_TOKEN=your_webhook_verify_token
   ```

## Troubleshooting

### Common Issues

1. **"Invalid phone number"**
   - Ensure phone number includes country code
   - Use format: `628123456789` (no +)

2. **"Access token expired"**
   - Generate a new access token
   - For production, use permanent token

3. **"Phone number not verified"**
   - Complete phone number verification in Meta Business Manager
   - Ensure your WhatsApp Business account is active

### Error Codes

- `100`: Invalid parameter
- `190`: Access token expired
- `368`: Temporary block (rate limiting)
- `131026`: Message undeliverable

## Rate Limits

Meta WhatsApp Business API has rate limits:
- **Tier 1**: 1,000 messages per day
- **Tier 2**: 10,000 messages per day
- **Tier 3**: 100,000 messages per day

## Cost

- First 1,000 conversations per month: Free
- After that: $0.005 - $0.05 per conversation (varies by country)

## Support

- [Meta WhatsApp Business API Documentation](https://developers.facebook.com/docs/whatsapp)
- [Meta Business Support](https://business.facebook.com/help)
