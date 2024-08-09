import 'package:json_annotation/json_annotation.dart';
import 'package:tutor/models/backend/chat_messages/get_chat_messages_response.dart';

part 'answer_to_messages_response.g.dart';

@JsonSerializable()
class AnswerToMessagesResponse {
  @JsonKey(name: 'msg')
  final ChatMessage Message;

  AnswerToMessagesResponse({
    required this.Message,
  });

  factory AnswerToMessagesResponse.fromJson(Map<String, dynamic> json) => _$AnswerToMessagesResponseFromJson(json);
  Map<String, dynamic> toJson() => _$AnswerToMessagesResponseToJson(this);
}
