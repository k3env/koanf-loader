package extensions

type TelegramConfig struct {
	Token      string `koanf:"token"`
	WebhookURL string `koanf:"webhookUrl"`
}
