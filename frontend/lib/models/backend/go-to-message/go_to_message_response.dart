import 'package:json_annotation/json_annotation.dart';
import 'package:tutor/models/backend/chat_messages/get_chat_messages_response.dart';

part 'go_to_message_response.g.dart';

@JsonSerializable()
class GoToMessageResponse {
  @JsonKey(name: 'msg')
  final ChatMessage Message;

  GoToMessageResponse({
    required this.Message,
  });

  factory GoToMessageResponse.fromJson(Map<String, dynamic> json) => _$GoToMessageResponseFromJson(json);
  Map<String, dynamic> toJson() => _$GoToMessageResponseToJson(this);
}