import 'package:json_annotation/json_annotation.dart';

part 'get_chats_response.g.dart';

@JsonSerializable()
class GetChatsResponse {
  @JsonKey(name: 'chats', defaultValue: <Chat>[])
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

  @JsonKey(name: 'typ', defaultValue: 0)
  final int ChatType;

  @JsonKey(name: 'cur_m', defaultValue: 0)
  final int CurrentMessageIDx;

  Chat({
    required this.ChatId,
    required this.Timestamp,
    required this.Title,
    required this.ChatType,
    required this.CurrentMessageIDx,
  });

  factory Chat.fromJson(Map<String, dynamic> json) => _$ChatFromJson(json);
  Map<String, dynamic> toJson() => _$ChatToJson(this);
}