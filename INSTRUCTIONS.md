# How to Run WhatsApp Service

## Step 1: Install Dependencies
```bash
cd whatsapp-service
npm install
```

## Step 2: Start the Service
```bash
npm start
```

Or on Windows:
```cmd
npm start
```

## Step 3: Scan QR Code (IMPORTANT!)

**DO NOT use your phone's camera app!** You must scan the QR code **directly from WhatsApp app**.

### Correct Steps:
1. **Open WhatsApp** on your phone (not camera app)
2. Tap **Settings** (gear icon, bottom right on Android, or top right on iOS)
3. Tap **Linked Devices** (or "Linked Devices" / "Paired Devices")
4. Tap **"Link a Device"** or **"+"** button
5. **Point your phone camera at the QR code** shown in the terminal
6. Wait for connection confirmation

### What NOT to do:
- ❌ Don't use your phone's Camera app
- ❌ Don't take a screenshot and scan from gallery
- ❌ Don't use a third-party QR scanner app

### Why it redirects to camera app:
If scanning redirects to camera app, it means you're using the wrong method. **You MUST scan from within WhatsApp app itself.**

## Step 4: Verify Connection

After scanning, you should see:
```
✅ [WA] Client is ready! WhatsApp connected successfully!
```

## Step 5: Test the Service

Check if service is ready:
```bash
curl http://localhost:3002/status
```

Or in browser: http://localhost:3002/status

## Troubleshooting

### "Client is not ready" after restart:

1. **Check if session folder exists:**
   ```bash
   ls whatsapp-session
   ```
   If it exists, the session should be preserved.

2. **If still not working, delete session and rescan:**
   ```bash
   rm -rf whatsapp-session
   npm start
   ```
   Then scan QR code again.

3. **Check logs** for any error messages

### Session not persisting:

- Make sure `whatsapp-session` folder has write permissions
- Don't delete the `whatsapp-session` folder after first scan
- The session persists between restarts automatically

### QR Code not showing:

- Make sure `qrcode-terminal` is installed: `npm install`
- Check terminal output for any errors
- Try restarting the service






