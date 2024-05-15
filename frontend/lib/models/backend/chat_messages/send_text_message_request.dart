import 'package:json_annotation/json_annotation.dart';

part 'send_text_message_request.g.dart';

@JsonSerializable()
class SendTextMessageRequest {
  @JsonKey(name: 'chatId')
  final String ChatId;

  @JsonKey(name: 'text')
  final String Text;

  @JsonKey(name: 'timestamp')
  final int Timestamp;

  SendTextMessageRequest({
    required this.ChatId,
    required this.Text,
    required this.Timestamp,
  });

  factory SendTextMessageRequest.fromJson(Map<String, dynamic> json) => _$SendTextMessageRequestFromJson(json);
  Map<String, dynamic> toJson() => _$SendTextMessageRequestToJson(this);
}