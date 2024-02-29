import 'package:json_annotation/json_annotation.dart';

part 'chat_message.g.dart';

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