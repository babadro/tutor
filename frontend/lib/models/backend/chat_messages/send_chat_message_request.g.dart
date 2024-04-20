// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'send_chat_message_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SendChatMessageRequest _$SendChatMessageRequestFromJson(
        Map<String, dynamic> json) =>
    SendChatMessageRequest(
      ChatId: json['chatId'] as String,
      Text: json['text'] as String,
      Timestamp: json['timestamp'] as int,
    );

Map<String, dynamic> _$SendChatMessageRequestToJson(
        SendChatMessageRequest instance) =>
    <String, dynamic>{
      'chatId': instance.ChatId,
      'text': instance.Text,
      'timestamp': instance.Timestamp,
    };
