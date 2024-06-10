import 'package:json_annotation/json_annotation.dart';

import 'package:tutor/models/backend/chats/get_chats_response.dart';

part 'send_text_message_response.g.dart';

@JsonSerializable()
class SendTextMessageResponse {
  @JsonKey(name: 'reply')
  final String Reply;

  @JsonKey(name: 'timestamp')
  final int Timestamp;

  @JsonKey(name: 'chat')
  final Chat? CreatedChat;

  SendTextMessageResponse({
    required this.Reply,
    required this.Timestamp,
    this.CreatedChat,
  });

  factory SendTextMessageResponse.fromJson(Map<String, dynamic> json) => _$SendTextMessageResponseFromJson(json);
  Map<String, dynamic> toJson() => _$SendTextMessageResponseToJson(this);
}