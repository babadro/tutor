// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'send_text_message_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SendTextMessageResponse _$SendTextMessageResponseFromJson(
        Map<String, dynamic> json) =>
    SendTextMessageResponse(
      Reply: json['reply'] as String,
      Timestamp: (json['timestamp'] as num).toInt(),
      CreatedChat: json['chat'] == null
          ? null
          : Chat.fromJson(json['chat'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$SendTextMessageResponseToJson(
        SendTextMessageResponse instance) =>
    <String, dynamic>{
      'reply': instance.Reply,
      'timestamp': instance.Timestamp,
      'chat': instance.CreatedChat,
    };
