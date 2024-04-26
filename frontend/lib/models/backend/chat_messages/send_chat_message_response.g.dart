// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'send_chat_message_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SendChatMessageResponse _$SendChatMessageResponseFromJson(
        Map<String, dynamic> json) =>
    SendChatMessageResponse(
      Reply: json['reply'] as String,
      Timestamp: json['timestamp'] as int,
      CreatedChat: json['chat'] == null
          ? null
          : Chat.fromJson(json['chat'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$SendChatMessageResponseToJson(
        SendChatMessageResponse instance) =>
    <String, dynamic>{
      'reply': instance.Reply,
      'timestamp': instance.Timestamp,
      'chat': instance.CreatedChat,
    };
