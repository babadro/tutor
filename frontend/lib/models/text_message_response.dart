import 'package:json_annotation/json_annotation.dart';

part 'text_message_response.g.dart';

@JsonSerializable()
class TextMessageResponse {
  @JsonKey(name: 'reply')
  final String Reply;

  @JsonKey(name: 'timestamp')
  final int Timestamp;

  TextMessageResponse({
    required this.Reply,
    required this.Timestamp,
  });

  factory TextMessageResponse.fromJson(Map<String, dynamic> json) => _$TextMessageResponseFromJson(json);
  Map<String, dynamic> toJson() => _$TextMessageResponseToJson(this);
}