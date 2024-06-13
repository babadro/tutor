// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'create_chat_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

CreateChatResponse _$CreateChatResponseFromJson(Map<String, dynamic> json) =>
    CreateChatResponse(
      CreatedChat: json['chat'] == null
          ? null
          : Chat.fromJson(json['chat'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$CreateChatResponseToJson(CreateChatResponse instance) =>
    <String, dynamic>{
      'chat': instance.CreatedChat,
    };
