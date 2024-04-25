import 'package:json_annotation/json_annotation.dart';

part 'send_chat_message_response.g.dart';

@JsonSerializable()
class SendChatMessageResponse {
  @JsonKey(name: 'reply')
  final String Reply;

  @JsonKey(name: 'timestamp')
  final int Timestamp;

  @JsonKey(name: 'chatId', defaultValue: '')
  final String ChatId;

  SendChatMessageResponse({
    required this.Reply,
    required this.Timestamp,
    required this.ChatId,
  });

  factory SendChatMessageResponse.fromJson(Map<String, dynamic> json) => _$SendChatMessageResponseFromJson(json);
  Map<String, dynamic> toJson() => _$SendChatMessageResponseToJson(this);
}