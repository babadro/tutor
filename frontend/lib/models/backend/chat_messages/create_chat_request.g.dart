// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'create_chat_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

CreateChatRequest _$CreateChatRequestFromJson(Map<String, dynamic> json) =>
    CreateChatRequest(
      ChatType: (json['typ'] as num).toInt(),
      Timestamp: (json['time'] as num).toInt(),
    );

Map<String, dynamic> _$CreateChatRequestToJson(CreateChatRequest instance) =>
    <String, dynamic>{
      'typ': instance.ChatType,
      'time': instance.Timestamp,
    };
