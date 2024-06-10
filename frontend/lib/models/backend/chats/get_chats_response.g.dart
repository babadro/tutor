// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'get_chats_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

GetChatsResponse _$GetChatsResponseFromJson(Map<String, dynamic> json) =>
    GetChatsResponse(
      Chats: (json['chats'] as List<dynamic>)
          .map((e) => Chat.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$GetChatsResponseToJson(GetChatsResponse instance) =>
    <String, dynamic>{
      'chats': instance.Chats,
    };

Chat _$ChatFromJson(Map<String, dynamic> json) => Chat(
      ChatId: json['chatId'] as String,
      Timestamp: (json['time'] as num?)?.toInt() ?? 0,
      Title: json['title'] as String? ?? '',
    );

Map<String, dynamic> _$ChatToJson(Chat instance) => <String, dynamic>{
      'chatId': instance.ChatId,
      'time': instance.Timestamp,
      'title': instance.Title,
    };
