import 'package:json_annotation/json_annotation.dart';

import 'package:tutor/models/backend/chats/get_chats_response.dart';

part 'create_chat_response.g.dart';

@JsonSerializable()
class CreateChatResponse {
  @JsonKey(name: 'chat')
  final Chat? CreatedChat;

  CreateChatResponse({
    this.CreatedChat,
  });

  factory CreateChatResponse.fromJson(Map<String, dynamic> json) => _$CreateChatResponseFromJson(json);
  Map<String, dynamic> toJson() => _$CreateChatResponseToJson(this);
}