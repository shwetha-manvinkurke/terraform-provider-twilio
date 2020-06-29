output "service" {
  description = "The Generated Chat Service"
  value       = twilio_chat_service.service
}

output "channel" {
  description = "The Generated Channel Channel"
  value       = twilio_chat_channel.channel
}
