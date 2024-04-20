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
    );

Map<String, dynamic> _$SendChatMessageResponseToJson(
        SendChatMessageResponse instance) =>
    <String, dynamic>{
      'reply': instance.Reply,
      'timestamp': instance.Timestamp,
    };
