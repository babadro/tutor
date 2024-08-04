import 'package:json_annotation/json_annotation.dart';

part 'go_to_message_request.g.dart';

@JsonSerializable()
class GoToMessageRequest {
  @JsonKey(name: 'chatId')
  final String ChatId;

  @JsonKey(name: 'msgIdx')
  final int MessageIndex;

  GoToMessageRequest({
    required this.ChatId,
    required this.MessageIndex,
  });

  factory GoToMessageRequest.fromJson(Map<String, dynamic> json) => _$GoToMessageRequestFromJson(json);
  Map<String, dynamic> toJson() => _$GoToMessageRequestToJson(this);
}