class ChatMessage {
  final bool IsFromCurrentUser;
  final String Text;
  final int Timestamp;

  ChatMessage({
    required this.IsFromCurrentUser,
    required this.Text,
    required this.Timestamp,
  });
}