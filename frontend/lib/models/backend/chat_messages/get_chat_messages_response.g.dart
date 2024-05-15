// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'get_chat_messages_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

GetChatMessagesResponse _$GetChatMessagesResponseFromJson(
        Map<String, dynamic> json) =>
    GetChatMessagesResponse(
      Messages: (json['messages'] as List<dynamic>)
          .map((e) => ChatMessage.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$GetChatMessagesResponseToJson(
        GetChatMessagesResponse instance) =>
    <String, dynamic>{
      'messages': instance.Messages,
    };

ChatMessage _$ChatMessageFromJson(Map<String, dynamic> json) => ChatMessage(
      IsFromCurrentUser: json['curUsr'] as bool? ?? false,
      Text: json['text'] as String,
      Timestamp: (json['timestamp'] as num).toInt(),
      UserId: json['userId'] as String? ?? '',
    );

Map<String, dynamic> _$ChatMessageToJson(ChatMessage instance) =>
    <String, dynamic>{
      'curUsr': instance.IsFromCurrentUser,
      'text': instance.Text,
      'timestamp': instance.Timestamp,
      'userId': instance.UserId,
    };
