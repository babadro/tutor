// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'chat_message.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ChatMessage _$ChatMessageFromJson(Map<String, dynamic> json) => ChatMessage(
      IsFromAi: json['isFromAi'] as bool,
      Text: json['text'] as String,
      Timestamp: json['timestamp'] as int,
      UserId: json['userId'] as String,
    );

Map<String, dynamic> _$ChatMessageToJson(ChatMessage instance) =>
    <String, dynamic>{
      'isFromAi': instance.IsFromAi,
      'text': instance.Text,
      'timestamp': instance.Timestamp,
      'userId': instance.UserId,
    };
