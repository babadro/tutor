import 'package:json_annotation/json_annotation.dart';

part 'chat_message.g.dart';

@JsonSerializable()
class ChatMessage {
  @JsonKey(name: 'isFromAi')
  final bool IsFromAi;

  @JsonKey(name: 'text')
  final String Text;

  @JsonKey(name: 'timestamp')
  final int Timestamp;

  @JsonKey(name: 'userId')
  final String UserId;

  ChatMessage({
    required this.IsFromAi,
    required this.Text,
    required this.Timestamp,
    required this.UserId,
  });

  factory ChatMessage.fromJson(Map<String, dynamic> json) => _$ChatMessageFromJson(json);
  Map<String, dynamic> toJson() => _$ChatMessageToJson(this);
}