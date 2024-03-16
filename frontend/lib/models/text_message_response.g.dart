// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'text_message_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

TextMessageResponse _$TextMessageResponseFromJson(Map<String, dynamic> json) =>
    TextMessageResponse(
      Reply: json['reply'] as String,
      Timestamp: json['timestamp'] as int,
    );

Map<String, dynamic> _$TextMessageResponseToJson(
        TextMessageResponse instance) =>
    <String, dynamic>{
      'reply': instance.Reply,
      'timestamp': instance.Timestamp,
    };
