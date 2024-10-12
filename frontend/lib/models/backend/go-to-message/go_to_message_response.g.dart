// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'go_to_message_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

GoToMessageResponse _$GoToMessageResponseFromJson(Map<String, dynamic> json) =>
    GoToMessageResponse(
      Message: ChatMessage.fromJson(json['msg'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$GoToMessageResponseToJson(
        GoToMessageResponse instance) =>
    <String, dynamic>{
      'msg': instance.Message,
    };
