import 'package:json_annotation/json_annotation.dart';

part 'get_chat_messages_response.g.dart';

@JsonSerializable()
class GetChatMessagesResponse {
  @JsonKey(name: 'messages')
  final List<ChatMessage> Messages;

  GetChatMessagesResponse({
    required this.Messages,
  });

  factory GetChatMessagesResponse.fromJson(Map<String, dynamic> json) => _$GetChatMessagesResponseFromJson(json);
  Map<String, dynamic> toJson() => _$GetChatMessagesResponseToJson(this);
}

@JsonSerializable()
class ChatMessage {
  @JsonKey(name: 'curUsr', defaultValue: false)
  final bool IsFromCurrentUser;

  @JsonKey(name: 'text')
  final String Text;

  @JsonKey(name: 'timestamp')
  final int Timestamp;

  @JsonKey(name: 'userId', defaultValue: '')
  final String UserId;

  ChatMessage({
    required this.IsFromCurrentUser,
    required this.Text,
    required this.Timestamp,
    required this.UserId,
  });

  factory ChatMessage.fromJson(Map<String, dynamic> json) => _$ChatMessageFromJson(json);
  Map<String, dynamic> toJson() => _$ChatMessageToJson(this);
}