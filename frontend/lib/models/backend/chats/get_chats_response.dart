import 'package:json_annotation/json_annotation.dart';

part 'get_chats_response.g.dart';

@JsonSerializable()
class GetChatsResponse {
  @JsonKey(name: 'chats')
  final List<Chat> Chats;

  GetChatsResponse({
    required this.Chats,
  });

  factory GetChatsResponse.fromJson(Map<String, dynamic> json) => _$GetChatsResponseFromJson(json);
  Map<String, dynamic> toJson() => _$GetChatsResponseToJson(this);
}

@JsonSerializable()
class Chat {
  @JsonKey(name: 'chatId')
  final String ChatId;

  @JsonKey(name: 'time', defaultValue: 0)
  final int Timestamp;

  @JsonKey(name: 'title', defaultValue: '')
  final String Title;

  Chat({
    required this.ChatId,
    required this.Timestamp,
    required this.Title,
  });

  factory Chat.fromJson(Map<String, dynamic> json) => _$ChatFromJson(json);
  Map<String, dynamic> toJson() => _$ChatToJson(this);
}