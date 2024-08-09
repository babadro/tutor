import 'package:json_annotation/json_annotation.dart';

part 'answer_to_messages_request.g.dart';

@JsonSerializable()
class AnswerToMessagesRequest {
  @JsonKey(name: 'chatId')
  final String ChatId;

  AnswerToMessagesRequest({
    required this.ChatId,
  });

  factory AnswerToMessagesRequest.fromJson(Map<String, dynamic> json) =>
      _$AnswerToMessagesRequestFromJson(json);
  Map<String, dynamic> toJson() => _$AnswerToMessagesRequestToJson(this);
}