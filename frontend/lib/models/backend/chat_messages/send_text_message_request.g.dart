// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'send_text_message_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SendTextMessageRequest _$SendTextMessageRequestFromJson(
        Map<String, dynamic> json) =>
    SendTextMessageRequest(
      ChatId: json['chatId'] as String,
      Text: json['text'] as String,
      Timestamp: (json['timestamp'] as num).toInt(),
    );

Map<String, dynamic> _$SendTextMessageRequestToJson(
        SendTextMessageRequest instance) =>
    <String, dynamic>{
      'chatId': instance.ChatId,
      'text': instance.Text,
      'timestamp': instance.Timestamp,
    };
