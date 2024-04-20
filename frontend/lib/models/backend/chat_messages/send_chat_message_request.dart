import 'package:json_annotation/json_annotation.dart';

part 'send_chat_message_request.g.dart';

@JsonSerializable()
class SendChatMessageRequest {
  @JsonKey(name: 'chatId')
  final String ChatId;

  @JsonKey(name: 'text')
  final String Text;

  @JsonKey(name: 'timestamp')
  final int Timestamp;

  SendChatMessageRequest({
    required this.ChatId,
    required this.Text,
    required this.Timestamp,
  });

  factory SendChatMessageRequest.fromJson(Map<String, dynamic> json) => _$SendChatMessageRequestFromJson(json);
  Map<String, dynamic> toJson() => _$SendChatMessageRequestToJson(this);
}