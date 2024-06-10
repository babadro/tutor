class ChatMessage {
  final bool IsFromCurrentUser;
  final String Text;
  final int Timestamp;
  final String AudioUrl;

  ChatMessage({
    required this.IsFromCurrentUser,
    required this.Text,
    required this.Timestamp,
    this.AudioUrl = '',
  });
}