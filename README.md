

## To host on Google Cloud Functions
(zipped)
- function.go
- go.mod
- go.sum
- services/

## Setting Webhook
TELEGRAM_TOKEN=""
CLOUD_FUNCTION_URL=""

curl --data "url=$CLOUD_FUNCTION_URL" https://api.telegram.org/bot$TELEGRAM_TOKEN/SetWebhook