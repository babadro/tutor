import 'package:json_annotation/json_annotation.dart';

part 'create_chat_request.g.dart';

@JsonSerializable()
class CreateChatRequest {
  @JsonKey(name: 'typ')
  final int ChatType;

  @JsonKey(name: 'time')
  final int Timestamp;

  CreateChatRequest({
    required this.ChatType,
    required this.Timestamp,
  });

  factory CreateChatRequest.fromJson(Map<String, dynamic> json) => _$CreateChatRequestFromJson(json);
  Map<String, dynamic> toJson() => _$CreateChatRequestToJson(this);
}
