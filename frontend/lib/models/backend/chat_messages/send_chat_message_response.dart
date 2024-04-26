import 'package:json_annotation/json_annotation.dart';

import 'package:tutor/models/backend/chats/get_chats_response.dart';

part 'send_chat_message_response.g.dart';

@JsonSerializable()
class SendChatMessageResponse {
  @JsonKey(name: 'reply')
  final String Reply;

  @JsonKey(name: 'timestamp')
  final int Timestamp;

  @JsonKey(name: 'chat')
  final Chat? CreatedChat;

  SendChatMessageResponse({
    required this.Reply,
    required this.Timestamp,
    this.CreatedChat,
  });

  factory SendChatMessageResponse.fromJson(Map<String, dynamic> json) => _$SendChatMessageResponseFromJson(json);
  Map<String, dynamic> toJson() => _$SendChatMessageResponseToJson(this);
}