// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'go_to_message_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

GoToMessageRequest _$GoToMessageRequestFromJson(Map<String, dynamic> json) =>
    GoToMessageRequest(
      ChatId: json['chatId'] as String,
      MessageIndex: (json['msgIdx'] as num).toInt(),
    );

Map<String, dynamic> _$GoToMessageRequestToJson(GoToMessageRequest instance) =>
    <String, dynamic>{
      'chatId': instance.ChatId,
      'msgIdx': instance.MessageIndex,
    };
