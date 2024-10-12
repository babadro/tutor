// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'answer_to_messages_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

AnswerToMessagesResponse _$AnswerToMessagesResponseFromJson(
        Map<String, dynamic> json) =>
    AnswerToMessagesResponse(
      Message: ChatMessage.fromJson(json['msg'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$AnswerToMessagesResponseToJson(
        AnswerToMessagesResponse instance) =>
    <String, dynamic>{
      'msg': instance.Message,
    };
